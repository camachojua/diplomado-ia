# Jesús Alberto Aguirre Caro

# Objectives

Compare a set of models based on the AlexNet architecture: a pretrained AlexNet, a pytorch trained AlexNet and a tensorflow trained AlexNet. (Will be done in Google Collab)

# Pretrained AlexNet

AlexNet was acquired from torchvision as a pretrained model. As well as the CIFAR10 dataset.

A few adjustments were made to make AlexNet compatible with the quantity of classes from the dataset, as seen in class.

Looping over the dataset a couple of times we can start the validation.

With CIFAR10 10 thousand images, AlexNet achieved 82.55% accuracy. Impressive for an off-the-shelf model that spit out all this in minutes.

# Pytorch AlexNet from scratch

I defined a model using the architecture for AlexNet which include the following layers:

-2d convolution (11 size kernel)
-2d max polling (3 size kernel)
-2d convolution (5 size kernel)
-2d max polling (3 size kernel)
-2d convolution (3 size kernel)
-2d convolution (3 size kernel)
-2d convolution (3 size kernel)
-2d max polling (3 size kernel)
-dense (6400 neurons)
-dense (4096 neurons)
-dense (4096 neurons to a final 10 outputs)

After 10 epochs a model with AlexNet's architecture achieved 80.91% accuracy. Still impressive considering it took almost an hour to run in collab.

# Tensorflow AlexNet from scratch

For the tensorflow version I found that Keras had its own way to download CIFAR10 and used its API to do so. The data came in 32x32 pixel size and I did not resize to the 224x224 size we have been using out of curiosity. Will the model have similar reults? Will it take less time? Would the same architecture work?

First of all, the architecture needed adjusting. By the 5fth comvolution layer, the data it received was too little to use a 3x3 kernel so I reduced it to a 1x1 layer to 'conserve' the outline of the archjitecture, instead of removing it.

After 10 epochs the model achieved 71.89% accuracy and took only 3 minutes to do so.

# Conlusions

The pretrained AlexNet was the undisputable victor. Which makes sense, presumably the time and resources it took to train this model was far greater than those I invested in this work.

But it serves as an excellent exercise and a confirmation that a good architecture will yield good results in general.

For an hour of training in pytorch it achieved 80.91% accuracy.
For 3 minutes and a reduced image size in tensorflow it achieved 71.89% accuracy.
