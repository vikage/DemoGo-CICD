# Install docker
```sudo apt-get install docker.io docker-compose```
# Deploy
- Step 1: Update .env.prod
- Step 2: Create alias (if not exist)
```ln -s .env.prod .env```
- Step 2: Run docker container
```make deploy```

# Gen mocks when want to add new mock
- Step 1: Install docker if not installed
- Step 2: Change directory to your project
- Step 3: Run script
```make gen```

# Run test
```make test```