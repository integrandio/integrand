FROM 3.12-alpine

WORKDIR /build

COPY integrand-py /build/integrand-py

WORKDIR /build/integrand-py
RUN pip install requirements.txt
ENV PYTHONPATH=src
CMD ["pytest"]