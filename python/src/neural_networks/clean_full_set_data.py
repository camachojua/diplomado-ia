def cleanFullSetData( dataDir):
    for image in sorted(glob.glob(dataDir)):
        try:
            # print(image)
            imgName = str(image)
            img = read_file( imgName )
            # img = decode_image(img)
            d_img = tf.image.decode_jpeg( img )
            # print( imgName )
            if d_img.shape[2] != 3:
                # os.remove(image)
                imgName.unlink()
                os.remove( imgName )
        except Exception as e:
            print(" bad file: ", imgName)
