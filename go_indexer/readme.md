docker run -v indexer:/data -e ZINC_DATA_PATH="/data" -p 4080:4080 \                                             
    -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 \
    --name zincsearch public.ecr.aws/zinclabs/zincsearch:latest
