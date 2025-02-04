# Dockerfile of image B
FROM a

RUN echo "this makes layer 4" > /four
RUN echo "this makes layer 5" > /five
