import random
import string
import re

import tensorflow.data as tf_data
import tensorflow.strings as tf_strings
from tensorflow.keras.callbacks import CSVLogger, ModelCheckpoint

import keras
import keras_nlp
from keras import layers
from keras import ops
from keras.layers import TextVectorization

"""
## Data
"""
text_file = "../corpus/spa.txt"


"""
## Parsing the data

Each line contains a French sentence and its corresponding Spanish sentence.
The French sentence is the *source sequence* and Spanish one is the *target sequence*.
We prepend the token `"[start]"` and we append the token `"[end]"` to the Spanish sentence.
"""

random.seed(10)

with open(text_file) as f:
    lines = f.read().split("\n")[:-1]
text_pairs = []
for line in lines:
    fre, spa = line.split("\t")
    spa = "[start] " + spa + " [end]"
    text_pairs.append((fre, spa))

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
data (one for French and one for Spanish),
that is to say, to turn the original strings into integer sequences
where each integer represents the index of a word in a vocabulary.

The French layer will use the default string standardization (strip punctuation characters)
and splitting scheme (split on whitespace), while
the Spanish layer will use a custom standardization, where we add the character
`"¿"` to the set of punctuation characters to be stripped.
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


fre_vectorization = TextVectorization(
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
train_fre_texts = [pair[0] for pair in train_pairs]
train_spa_texts = [pair[1] for pair in train_pairs]
fre_vectorization.adapt(train_fre_texts)
spa_vectorization.adapt(train_spa_texts)

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


def format_dataset(fre, spa):
    fre = fre_vectorization(fre)
    spa = spa_vectorization(spa)
    return (
        {
            "encoder_inputs": fre,
            "decoder_inputs": spa[:, :-1],
        },
        spa[:, 1:],
    )


def make_dataset(pairs):
    fre_texts, spa_texts = zip(*pairs)
    fre_texts = list(fre_texts)
    spa_texts = list(spa_texts)
    dataset = tf_data.Dataset.from_tensor_slices((fre_texts, spa_texts))
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

"""
## Building the model

Our sequence-to-sequence Transformer consists of a `TransformerEncoder`
and a `TransformerDecoder` chained together. To make the model aware of word order,
we also use a `PositionalEmbedding` layer.

The source sequence will be pass to the `TransformerEncoder`,
which will produce a new representation of it.
This new representation will then be passed
to the `TransformerDecoder`, together with the target sequence so far (target words 0 to N).
The `TransformerDecoder` will then seek to predict the next words in the target sequence (N+1 and beyond).

A key detail that makes this possible is causal masking
(see method `get_causal_attention_mask()` on the `TransformerDecoder`).
The `TransformerDecoder` sees the entire sequences at once, and thus we must make
sure that it only uses information from target tokens 0 to N when predicting token N+1
(otherwise, it could use information from the future, which would
result in a model that cannot be used at inference time).
"""


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
    def __init__(self, sequence_length, vocab_size, embed_dim, **kwargs):
        super().__init__(**kwargs)
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


"""
Next, we assemble the end-to-end model.
"""

embed_dim = 2**8
latent_dim = 2**11
num_heads = 8

encoder_inputs = keras.Input(shape=(None,), dtype="int64", name="encoder_inputs")
x = PositionalEmbedding(sequence_length, vocab_size, embed_dim)(encoder_inputs)
encoder_outputs = TransformerEncoder(embed_dim, latent_dim, num_heads)(x)
encoder = keras.Model(encoder_inputs, encoder_outputs)

decoder_inputs = keras.Input(shape=(None,), dtype="int64", name="decoder_inputs")
encoded_seq_inputs = keras.Input(shape=(None, embed_dim), name="decoder_state_inputs")
x = PositionalEmbedding(sequence_length, vocab_size, embed_dim)(decoder_inputs)
x = TransformerDecoder(embed_dim, latent_dim, num_heads)([x, encoder_outputs])
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

transformer.summary()
# transformer.compile(
#   "rmsprop",
#   loss=keras.losses.SparseCategoricalCrossentropy(ignore_class=0),
#   metrics=["accuracy"],
# )
epochs = 5  # This should be at least 30 for convergence
transformer = keras.models.load_model(
    "../transformers/transformer_chollet_glove100d.checkpoint.keras",
    custom_objects={
        "encoder_inputs": encoder_inputs,
        "decoder_inputs": decoder_inputs,
        "TransformerEncoder": TransformerEncoder,
        "TransformerDecoder": TransformerDecoder,
        "PositionalEmbedding": PositionalEmbedding,
    },
)

# csv logger callback
csv_logger = CSVLogger("transformer_v4.training.log")

# model checkpoint callback
model_checkpoint = ModelCheckpoint(
    filepath="transformer_v4.checkpoint.keras",
    monitor="val_accuracy",
    mode="max",
    save_best_only=True,
)
# transformer.fit(
#    train_ds,
#    epochs=epochs,
#    validation_data=val_ds,
#    verbose=2,
#    callbacks=[csv_logger, model_checkpoint],
# )

"""
## Decoding test sentences

Finally, let's demonstrate how to translate brand new French sentences.
We simply feed into the model the vectorized French sentence
as well as the target token `"[start]"`, then we repeatedly generated the next token, until
we hit the token `"[end]"`.
"""

spa_vocab = spa_vectorization.get_vocabulary()
spa_index_lookup = dict(zip(range(len(spa_vocab)), spa_vocab))
max_decoded_sentence_length = 20


def decode_sequence(input_sentence):
    tokenized_input_sentence = fre_vectorization([input_sentence])
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


random.seed(10)
test_fre_texts = [pair[0] for pair in test_pairs]
for _ in range(10):
    input_sentence = random.choice(test_fre_texts)
    translated = decode_sequence(input_sentence)
    print(input_sentence)
    print(translated)
    print()


"""
## Quantitative evaluation (Rouge metrics)
"""
rouge_1 = keras_nlp.metrics.RougeN(order=1)
rouge_2 = keras_nlp.metrics.RougeN(order=2)
bleu_1 = keras_nlp.metrics.Bleu(max_order=1)
bleu_2 = keras_nlp.metrics.Bleu(max_order=2)

reference_sentences = []
translated_sentences = []
for test_pair in test_pairs[:30]:
    input_sentence = test_pair[0]
    reference_sentence = test_pair[1]

    translated_sentence = decode_sequence(input_sentence)
    translated_sentence = (
        translated_sentence.replace("[PAD]", "")
        .replace("[START]", "")
        .replace("[END]", "")
        .strip()
    )
    reference_sentences.append(reference_sentence)
    translated_sentences.append(translated_sentence)

rouge_1(reference_sentences, translated_sentences)
rouge_2(reference_sentences, translated_sentences)
bleu_1(reference_sentences, translated_sentences)
bleu_2(reference_sentences, translated_sentences)

print("ROUGE-1 Score: ", rouge_1.result())
print("ROUGE-2 Score: ", rouge_2.result())
print("BLEU-1 Score: ", bleu_1.result())
print("BLEU-2 Score: ", bleu_2.result())


# import matplotlib.pyplot as plt
#
#
# def plot_attention_head(in_tokens, translated_tokens, attention):
#    # The model didn't generate `` in the output. Skip it.
#    translated_tokens = translated_tokens[1:]
#
#    ax = plt.gca()
#    ax.matshow(attention)
#    ax.set_xticks(range(len(in_tokens)))
#   ax.set_yticks(range(len(translated_tokens)))
#
#    labels = [label.decode("utf-8") for label in in_tokens.numpy()]
#    ax.set_xticklabels(labels, rotation=90)
#
#    labels = [label.decode("utf-8") for label in translated_tokens.numpy()]
#    ax.set_yticklabels(labels)
#
#
# def plot_attention_weights(sentence):
#    tokenized_input_sentence = fre_vectorization([sentence])
#    translated_tokens = decode_sequence(sentence)
#
#    fig = plt.figure(figsize=(16, 8))
#
#    for h, head in enumerate(attention_scores):
#        ax = fig.add_subplot(2, 4, h + 1)
#
#        plot_attention_head(tokenized_input_sentence, translated_tokens, head)
#
#        ax.set_xlabel(f"Head {h+1}")
#
#    plt.tight_layout()
#    plt.show()
#    fig.savefig("temp.png", dpi=fig.dpi)


# plot_attention_weights("tom is looking for you")
#model = keras.Model(
#    inputs=transformer.input,
#    outputs=[
#        transformer.output,
#        transformer.get_layer("transformer_encoder").output,
#        transformer.get_layer("transformer_decoder").output,
#    ],
#)
##https://data-science-blog.com/blog/2021/04/07/multi-head-attention-mechanism/
## https://github.com/TaiToTo/Transformer_blog_codes/blob/main/show_self_attention_en_es.ipynb
#
#
#def calculate_attention_weights(input_sentence, model):
#    tokenized_input_sentence = fre_vectorization([input_sentence])
#    decoded_sentence = "[start]"
#
#    key_embeddings = []
#    query_embeddings = []
#    for i in range(max_decoded_sentence_length):
#        tokenized_target_sentence = spa_vectorization([decoded_sentence])[:, :-1]
#        predictions = model(
#            {
#                "encoder_inputs": tokenized_input_sentence,
#                "decoder_inputs": tokenized_target_sentence,
#            }
#        )
#
#        sampled_token_index = ops.convert_to_numpy(
#            ops.argmax(predictions[0][0, i, :])
#        ).item(0)
#        sampled_token = spa_index_lookup[sampled_token_index]
#        decoded_sentence += " " + sampled_token
#
#        print(predictions[1][0, i, :])
#        print(predictions[2][0, i, :])
#
#        if sampled_token == "[end]":
#            break
#    return decoded_sentence
#
#
#cucu = calculate_attention_weights("Tom is well aware of the problem.", model)
#print(cucu)
