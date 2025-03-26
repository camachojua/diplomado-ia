import os

os.environ["KERAS_BACKEND"] = "tensorflow"

import pathlib
import random
import string
import re
import numpy as np

import tensorflow.data as tf_data
import tensorflow.strings as tf_strings

import keras
from keras import layers
from keras import ops
from keras.layers import TextVectorization

import time

text_file = 'spa.txt'

"""
## Parsing the data

Each line contains an English sentence and its corresponding Spanish sentence.
The English sentence is the *source sequence* and Spanish one is the *target sequence*.
We prepend the token `"[start]"` and we append the token `"[end]"` to the Spanish sentence.
"""

with open(text_file) as f:
    lines = f.read().split("\n")[:-1]
text_pairs = []
for line in lines:
    eng, spa = line.split("\t")
    spa = "[start] " + spa + " [end]"
    text_pairs.append((eng, spa))

"""
Here's what our sentence pairs look like:
"""

for _ in range(5):
    print(random.choice(text_pairs))

"""
Now, let's split the sentence pairs into a training set, a validation set,
and a test set.
"""

random.shuffle(text_pairs)
num_val_samples = int(0.15 * len(text_pairs))
num_train_samples = len(text_pairs) - 2 * num_val_samples
train_pairs = text_pairs[:num_train_samples]
val_pairs = text_pairs[num_train_samples : num_train_samples + num_val_samples]
test_pairs = text_pairs[num_train_samples + num_val_samples :]

print(f"{len(text_pairs)} total pairs")
print(f"{len(train_pairs)} training pairs")
print(f"{len(val_pairs)} validation pairs")
print(f"{len(test_pairs)} test pairs")

"""
## Vectorizing the text data

We'll use two instances of the `TextVectorization` layer to vectorize the text
data (one for English and one for Spanish),
that is to say, to turn the original strings into integer sequences
where each integer represents the index of a word in a vocabulary.

The English layer will use the default string standardization (strip punctuation characters)
and splitting scheme (split on whitespace), while
the Spanish layer will use a custom standardization, where we add the character
`"¿"` to the set of punctuation characters to be stripped.

Note: in a production-grade machine translation model, I would not recommend
stripping the punctuation characters in either language. Instead, I would recommend turning
each punctuation character into its own token,
which you could achieve by providing a custom `split` function to the `TextVectorization` layer.
"""

strip_chars = string.punctuation + "¿"
strip_chars = strip_chars.replace("[", "")
strip_chars = strip_chars.replace("]", "")

vocab_size = 15000
sequence_length = 20
batch_size = 64

def custom_standardization(input_string):
    lowercase = tf_strings.lower(input_string)
    return tf_strings.regex_replace(lowercase, "[%s]" % re.escape(strip_chars), "")


eng_vectorization = TextVectorization(
    max_tokens=vocab_size,
    output_mode="int",
    output_sequence_length=sequence_length,
)
spa_vectorization = TextVectorization(
    max_tokens=vocab_size,
    output_mode="int",
    output_sequence_length=sequence_length + 1,
    standardize=custom_standardization,
)
train_eng_texts = [pair[0] for pair in train_pairs]
train_spa_texts = [pair[1] for pair in train_pairs]
eng_vectorization.adapt(train_eng_texts)
spa_vectorization.adapt(train_spa_texts)

"""
Codigo para agregar los word embeddings de standord
"""
glove_file = "glove.6B.100d.txt"

embedding_index = {}
with open(glove_file, encoding="utf8") as f:
    for line in f:
        values = line.split()
        word = values[0]
        coefs = np.asarray(values[1:], dtype="float32")
        embedding_index[word] = coefs

# Crear la matriz de embeddings para el vocabulario en inglés
eng_vocab = eng_vectorization.get_vocabulary()
embedding_dim = 100  # Dimension de los embeddings de GloVe que se usan
num_tokens = len(eng_vocab)
embedding_matrix = np.zeros((num_tokens, embedding_dim))
for i, word in enumerate(eng_vocab):
    embedding_vector = embedding_index.get(word)
    if embedding_vector is not None:
        embedding_matrix[i] = embedding_vector
    # Si la palabra no está en GloVe, se deja el vector de ceros

print("Matriz de embeddings creada con forma:", embedding_matrix.shape)

"""
Next, we'll format our datasets.

At each training step, the model will seek to predict target words N+1 (and beyond)
using the source sentence and the target words 0 to N.

As such, the training dataset will yield a tuple `(inputs, targets)`, where:

- `inputs` is a dictionary with the keys `encoder_inputs` and `decoder_inputs`.
`encoder_inputs` is the vectorized source sentence and `decoder_inputs` is the target sentence "so far",
that is to say, the words 0 to N used to predict word N+1 (and beyond) in the target sentence.
- `target` is the target sentence offset by one step:
it provides the next words in the target sentence -- what the model will try to predict.
"""

def format_dataset(eng, spa):
    eng = eng_vectorization(eng)
    spa = spa_vectorization(spa)
    return (
        {
            "encoder_inputs": eng,
            "decoder_inputs": spa[:, :-1],
        },
        spa[:, 1:],
    )


def make_dataset(pairs):
    eng_texts, spa_texts = zip(*pairs)
    eng_texts = list(eng_texts)
    spa_texts = list(spa_texts)
    dataset = tf_data.Dataset.from_tensor_slices((eng_texts, spa_texts))
    dataset = dataset.batch(batch_size)
    dataset = dataset.map(format_dataset)
    return dataset.cache().shuffle(2048).prefetch(16)


train_ds = make_dataset(train_pairs)
val_ds = make_dataset(val_pairs)

"""
Let's take a quick look at the sequence shapes
(we have batches of 64 pairs, and all sequences are 20 steps long):
"""

for inputs, targets in train_ds.take(1):
    print(f'inputs["encoder_inputs"].shape: {inputs["encoder_inputs"].shape}')
    print(f'inputs["decoder_inputs"].shape: {inputs["decoder_inputs"].shape}')
    print(f"targets.shape: {targets.shape}")

import keras.ops as ops

class TransformerEncoder(layers.Layer):
    def __init__(self, embed_dim, dense_dim, num_heads, **kwargs):
        super().__init__(**kwargs)
        self.embed_dim = embed_dim
        self.dense_dim = dense_dim
        self.num_heads = num_heads
        self.attention = layers.MultiHeadAttention(
            num_heads=num_heads, key_dim=embed_dim
        )
        self.dense_proj = keras.Sequential(
            [
                layers.Dense(dense_dim, activation="relu"),
                layers.Dense(embed_dim),
            ]
        )
        self.layernorm_1 = layers.LayerNormalization()
        self.layernorm_2 = layers.LayerNormalization()
        self.supports_masking = True

    def call(self, inputs, mask=None):
        if mask is not None:
            padding_mask = ops.cast(mask[:, None, :], dtype="int32")
        else:
            padding_mask = None

        attention_output = self.attention(
            query=inputs, value=inputs, key=inputs, attention_mask=padding_mask
        )
        proj_input = self.layernorm_1(inputs + attention_output)
        proj_output = self.dense_proj(proj_input)
        return self.layernorm_2(proj_input + proj_output)

    def get_config(self):
        config = super().get_config()
        config.update(
            {
                "embed_dim": self.embed_dim,
                "dense_dim": self.dense_dim,
                "num_heads": self.num_heads,
            }
        )
        return config


class PositionalEmbedding(layers.Layer):
    def __init__(self, sequence_length, vocab_size, embed_dim, embedding_matrix=None, trainable=True, **kwargs):
        super().__init__(**kwargs)
        if embedding_matrix is not None:
            
            self.token_embeddings = layers.Embedding(
                input_dim=vocab_size,
                output_dim=embed_dim,
                weights=[embedding_matrix],
                trainable=trainable
            )
        else:
            self.token_embeddings = layers.Embedding(
                input_dim=vocab_size, output_dim=embed_dim
            )
        self.position_embeddings = layers.Embedding(
            input_dim=sequence_length, output_dim=embed_dim
        )
        self.sequence_length = sequence_length
        self.vocab_size = vocab_size
        self.embed_dim = embed_dim

    def call(self, inputs):
        length = ops.shape(inputs)[-1]
        positions = ops.arange(0, length, 1)
        embedded_tokens = self.token_embeddings(inputs)
        embedded_positions = self.position_embeddings(positions)
        return embedded_tokens + embedded_positions

    def compute_mask(self, inputs, mask=None):
        return ops.not_equal(inputs, 0)

    def get_config(self):
        config = super().get_config()
        config.update(
            {
                "sequence_length": self.sequence_length,
                "vocab_size": self.vocab_size,
                "embed_dim": self.embed_dim,
            }
        )
        return config


class TransformerDecoder(layers.Layer):
    def __init__(self, embed_dim, latent_dim, num_heads, **kwargs):
        super().__init__(**kwargs)
        self.embed_dim = embed_dim
        self.latent_dim = latent_dim
        self.num_heads = num_heads
        self.attention_1 = layers.MultiHeadAttention(
            num_heads=num_heads, key_dim=embed_dim
        )
        self.attention_2 = layers.MultiHeadAttention(
            num_heads=num_heads, key_dim=embed_dim
        )
        self.dense_proj = keras.Sequential(
            [
                layers.Dense(latent_dim, activation="relu"),
                layers.Dense(embed_dim),
            ]
        )
        self.layernorm_1 = layers.LayerNormalization()
        self.layernorm_2 = layers.LayerNormalization()
        self.layernorm_3 = layers.LayerNormalization()
        self.supports_masking = True

    def call(self, inputs, mask=None):
        inputs, encoder_outputs = inputs
        causal_mask = self.get_causal_attention_mask(inputs)

        if mask is None:
            inputs_padding_mask, encoder_outputs_padding_mask = None, None
        else:
            inputs_padding_mask, encoder_outputs_padding_mask = mask

        attention_output_1 = self.attention_1(
            query=inputs,
            value=inputs,
            key=inputs,
            attention_mask=causal_mask,
            query_mask=inputs_padding_mask,
        )
        out_1 = self.layernorm_1(inputs + attention_output_1)

        attention_output_2 = self.attention_2(
            query=out_1,
            value=encoder_outputs,
            key=encoder_outputs,
            query_mask=inputs_padding_mask,
            key_mask=encoder_outputs_padding_mask,
        )
        out_2 = self.layernorm_2(out_1 + attention_output_2)

        proj_output = self.dense_proj(out_2)
        return self.layernorm_3(out_2 + proj_output)

    def get_causal_attention_mask(self, inputs):
        input_shape = ops.shape(inputs)
        batch_size, sequence_length = input_shape[0], input_shape[1]
        i = ops.arange(sequence_length)[:, None]
        j = ops.arange(sequence_length)
        mask = ops.cast(i >= j, dtype="int32")
        mask = ops.reshape(mask, (1, input_shape[1], input_shape[1]))
        mult = ops.concatenate(
            [ops.expand_dims(batch_size, -1), ops.convert_to_tensor([1, 1])],
            axis=0,
        )
        return ops.tile(mask, mult)
    
    def get_attention_scores(self, inputs, encoder_outputs):
        """
        Devuelve las puntuaciones de atención de la capa de atención cruzada.
        inputs: secuencia de entrada al decodificador (sin el último token)
        encoder_outputs: salida del encoder
        """
        causal_mask = self.get_causal_attention_mask(inputs)
        
        attention_scores = self.attention_2(
            query=inputs,
            value=encoder_outputs,
            key=encoder_outputs,
            attention_mask=causal_mask,
            return_attention_scores=True,
        )[1]
        return attention_scores

    def get_config(self):
        config = super().get_config()
        config.update(
            {
                "embed_dim": self.embed_dim,
                "latent_dim": self.latent_dim,
                "num_heads": self.num_heads,
            }
        )
        return config
    
eng_vocab = eng_vectorization.get_vocabulary()
eng_vocab_size = len(eng_vocab)

encoder_embed_dim = embedding_dim
latent_dim = 2048
num_heads = 8

encoder_inputs = keras.Input(shape=(None,), dtype="int64", name="encoder_inputs")

x = PositionalEmbedding(sequence_length, eng_vocab_size, encoder_embed_dim, 
                          embedding_matrix=embedding_matrix, trainable=False)(encoder_inputs)
encoder_outputs = TransformerEncoder(encoder_embed_dim, latent_dim, num_heads)(x)
encoder = keras.Model(encoder_inputs, encoder_outputs)

decoder_inputs = keras.Input(shape=(None,), dtype="int64", name="decoder_inputs")
encoded_seq_inputs = keras.Input(shape=(None, encoder_embed_dim), name="decoder_state_inputs")

x = PositionalEmbedding(sequence_length, vocab_size, encoder_embed_dim, name="decoder_embedding")(decoder_inputs)

x = TransformerDecoder(encoder_embed_dim, latent_dim, num_heads, name="transformer_decoder")([x, encoder_outputs])
x = layers.Dropout(0.5)(x)
decoder_outputs = layers.Dense(vocab_size, activation="softmax")(x)
decoder = keras.Model([decoder_inputs, encoded_seq_inputs], decoder_outputs)

transformer = keras.Model(
    {"encoder_inputs": encoder_inputs, "decoder_inputs": decoder_inputs},
    decoder_outputs,
    name="transformer",
)

"""
## Training our model

We'll use accuracy as a quick way to monitor training progress on the validation data.
Note that machine translation typically uses BLEU scores as well as other metrics, rather than accuracy.

Here we only train for 1 epoch, but to get the model to actually converge
you should train for at least 30 epochs.
"""
inicio = time.time()
epochs = 35  # This should be at least 30 for convergence

#Cambie el optimizador por adam
transformer.summary()
transformer.compile(
    "adam",
    loss=keras.losses.SparseCategoricalCrossentropy(ignore_class=0),
    metrics=["accuracy"],
)
transformer.fit(train_ds, epochs=epochs, validation_data=val_ds)
final = time.time()

"""
## Decoding test sentences

Finally, let's demonstrate how to translate brand new English sentences.
We simply feed into the model the vectorized English sentence
as well as the target token `"[start]"`, then we repeatedly generated the next token, until
we hit the token `"[end]"`.
"""

spa_vocab = spa_vectorization.get_vocabulary()
spa_index_lookup = dict(zip(range(len(spa_vocab)), spa_vocab))
max_decoded_sentence_length = 20


def decode_sequence(input_sentence):
    tokenized_input_sentence = eng_vectorization([input_sentence])
    decoded_sentence = "[start]"
    for i in range(max_decoded_sentence_length):
        tokenized_target_sentence = spa_vectorization([decoded_sentence])[:, :-1]
        predictions = transformer(
            {
                "encoder_inputs": tokenized_input_sentence,
                "decoder_inputs": tokenized_target_sentence,
            }
        )

        # ops.argmax(predictions[0, i, :]) is not a concrete value for jax here
        sampled_token_index = ops.convert_to_numpy(
            ops.argmax(predictions[0, i, :])
        ).item(0)
        sampled_token = spa_index_lookup[sampled_token_index]
        decoded_sentence += " " + sampled_token

        if sampled_token == "[end]":
            break
    return decoded_sentence


test_eng_texts = [pair[0] for pair in test_pairs]
for _ in range(10):
    input_sentence = random.choice(test_eng_texts)
    translated = decode_sequence(input_sentence)
    print(input_sentence)
    print(translated)

def decode_sequence_with_attention(input_sentence):
    
    tokenized_input_sentence = eng_vectorization([input_sentence])
    encoder_out = encoder.predict(tokenized_input_sentence)
    
    decoded_sentence = "[start]"
    for i in range(max_decoded_sentence_length):
        tokenized_target_sentence = spa_vectorization([decoded_sentence])[:, :-1]
        predictions = transformer(
            {
                "encoder_inputs": tokenized_input_sentence,
                "decoder_inputs": tokenized_target_sentence,
            }
        )
        
        sampled_token_index = ops.convert_to_numpy(ops.argmax(predictions[0, i, :])).item(0)
        sampled_token = spa_index_lookup[sampled_token_index]
        decoded_sentence += " " + sampled_token
        if sampled_token == "[end]":
            break
    
    decoder_embedding_layer = transformer.get_layer("decoder_embedding")
    tokenized_target_sentence = spa_vectorization([decoded_sentence])[:, :-1]
    query_embeddings = decoder_embedding_layer(tokenized_target_sentence)
    
    transformer_decoder_layer = transformer.get_layer("transformer_decoder")
    attention_scores = transformer_decoder_layer.get_attention_scores(query_embeddings, encoder_out)
    
    return decoded_sentence, encoder_out, query_embeddings, attention_scores

#GRAFICAR EL PROCESO DE ATENCION

import math
import matplotlib.pyplot as plt
import seaborn as sns

def plot_attention_head(input_tokens, output_tokens, attention, head_idx=0):
    """
    Grafica un heatmap de la matriz de atención para un head específico,
    colocando los tokens de entrada (input) en el eje Y y los tokens de salida (output) en el eje X.
    Además, rota los tokens de salida para que se muestren en vertical.
    """
    
    attn = attention[head_idx].numpy()
    attn = attn[:len(output_tokens), :len(input_tokens)]
    attn = attn.T

    fig, ax = plt.subplots(figsize=(8, 6))
    sns.heatmap(
        attn,
        xticklabels=output_tokens,
        yticklabels=input_tokens,
        cmap='viridis',
        ax=ax
    )

    # Eje Y = Input
    ax.set_ylabel("Input")

    # Eje X = Output en la parte superior
    ax.set_xlabel("Output", labelpad=10)
    ax.xaxis.set_label_position('top')  # Etiqueta del eje X arriba
    ax.xaxis.tick_top()                # Ticks en la parte superior

    # Rotar los tokens de salida (eje X) para que queden verticales
    plt.setp(
        ax.get_xticklabels(),
        rotation=90,
        ha='center',
        va='bottom'
    )

    # Quitamos el título y lo colocamos en la parte inferior
    ax.set_title("")
    ax.text(
        0.5, -0.12,
        f"Head {head_idx+1}",
        ha='center', va='center',
        transform=ax.transAxes,
        fontsize=12
    )
    # plt.show()


def plot_all_attention_heads(input_tokens, output_tokens, attention, save_path=None):
    """
    Grafica la atención de todos los heads en subplots y opcionalmente la guarda en un archivo.
    Eje Y = tokens de entrada, Eje X (arriba) = tokens de salida en vertical.
    """
    num_heads = attention.shape[0]
    ncols = 4
    nrows = math.ceil(num_heads / ncols)

    fig, axes = plt.subplots(nrows=nrows, ncols=ncols, figsize=(20, nrows * 5))
    axes = axes.flatten()
    
    for i in range(num_heads):
        attn = attention[i].numpy()
        attn = attn[:len(output_tokens), :len(input_tokens)]
        attn = attn.T

        ax = axes[i]
        sns.heatmap(
            attn,
            xticklabels=output_tokens,
            yticklabels=input_tokens,
            cmap='viridis',
            ax=ax
        )

        # Eje Y = Input
        ax.set_ylabel("Input")

        # Eje X = Output arriba
        ax.set_xlabel("Output", labelpad=10)
        ax.xaxis.set_label_position('top')
        ax.xaxis.tick_top()

        # Rotamos las etiquetas del eje X para que se muestren en vertical
        plt.setp(
            ax.get_xticklabels(),
            rotation=90,
            ha='center',
            va='bottom'
        )

        # Titulo en la parte inferior
        ax.set_title("")
        ax.text(
            0.5, -0.12,
            f"Head {i+1}",
            ha='center', va='center',
            transform=ax.transAxes,
            fontsize=12
        )
    
    for j in range(i+1, len(axes)):
        axes[j].axis('off')
    
    plt.tight_layout()

    if save_path is not None:
        fig.savefig(save_path)

    # plt.show()

def translate_and_plot_attention(input_sentence, save_path_image):
    """
    Decodifica la oración de entrada, obtiene los pesos de atención y grafica
    los heatmaps de todos los heads.
    """
    # Decodificamos y obtenemos la oración traducida, la salida del encoder y los pesos de atencion
    decoded_sentence, encoder_out, query_embeddings, attention_scores = decode_sequence_with_attention(input_sentence)

    # Obtenemos los tokens de entrada y salida para etiquetar el heatmap
    eng_vocab = eng_vectorization.get_vocabulary()
    spa_vocab = spa_vectorization.get_vocabulary()
    
    tokenized_input = eng_vectorization([input_sentence]).numpy()[0]
    input_tokens = [eng_vocab[i] for i in tokenized_input if i != 0]

    tokenized_output = spa_vectorization([decoded_sentence]).numpy()[0]
    output_tokens = [spa_vocab[i] for i in tokenized_output if i != 0 and spa_vocab[i] != "[end]"]
    
    attention_scores = attention_scores[0]

    # Graficamos todos los heads
    print("Input sentence:", input_sentence)
    print("Translated sentence:", decoded_sentence)
    plot_all_attention_heads(input_tokens, output_tokens, attention_scores, save_path=save_path_image)

test_eng_texts = [pair[0] for pair in test_pairs]
example_sentence = random.choice(test_eng_texts)
translate_and_plot_attention(example_sentence, "all_heads_attention.png")

import pickle

# Guardar vocabulario en inglés y español
eng_vocab = eng_vectorization.get_vocabulary()
spa_vocab = spa_vectorization.get_vocabulary()

with open('eng_vocab.pkl', 'wb') as f:
    pickle.dump(eng_vocab, f)
with open('spa_vocab.pkl', 'wb') as f:
    pickle.dump(spa_vocab, f)

# Guardar diccionario de índices del vocabulario español (usado para la decodificación)
with open('spa_index_lookup.pkl', 'wb') as f:
    pickle.dump(spa_index_lookup, f)

# Guardar la matriz de embeddings
np.save('embedding_matrix.npy', embedding_matrix)

# Guardar el modelo
transformer.save('my_transformer.keras')
print(f"Tiempo de ejecución: {final - inicio:.6f} segundos")