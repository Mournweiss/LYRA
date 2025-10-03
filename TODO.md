<div align="center">

# TODO

Development roadmap

</div>

## Bug Fixes

-   [ ] Fix hardcoded service ports in compose.yml

## Refinements

-   [x] Add centralized error handling in the api-gateway
-   [x] Add centralized error handling in the telegram-bot
-   [x] Add centralized error handling in the whisper-service
-   [x] Add centralized config handling in the api-gateway
-   [x] Add centralized config handling in the telegram-bot
-   [x] Add centralized config handling in the whisper-service
-   [ ] Replace the file processing stub by accessing whisper.cpp in the whisper-service
-   [ ] Add validation of uploaded files
-   [x] Add direct transmission of environment variables to services in compose.yml
-   [ ] Integrate Redis to manage task statuses and queues
-   [x] Add default settings in config scripts
