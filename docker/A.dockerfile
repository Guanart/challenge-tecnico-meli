# Dockerfile of image A
FROM alpine
# layers 1 2 and 3
RUN echo "this makes layer 1" > /one
RUN echo "this makes layer 2" > /two
RUN echo "this makes layer 3" > /three
