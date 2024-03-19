# Frontend

Run the following command to start the frontend:

```bash
    cd data-visualizer
    npm i
    npm run dev
```

# Backend

Run the following command to start the backend:

```bash
    cd go_indexer_server
    ./com.nestorsgarzonc.indexer-server
```

# Indexer

Download the dataset from the following link:

```bash
    cd go_indexer_server
    wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz
    tar -xvzf enron_mail_20110402.tgz
```

Run the following command to start the indexer:

```bash
    cd go_indexer_server
    ./com.nestorsgarzonc.go-indexer
```

Indexer performance:

<img src="go_indexer/pprof001.svg" alt="drawing" width="500"/>

# ZincDB

Run the following command to start the ZincDB:

```bash
   docker run -v indexer:/data -e ZINC_DATA_PATH="/data" -p 4080:4080 \
    -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 \
    --name zincsearch public.ecr.aws/zinclabs/zincsearch:latest
```
