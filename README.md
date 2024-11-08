# Image Downloader
 
client application to search for images and storing them locally.

### How to run
Please first set the `API_KEY` in config.yaml. You can get from [PEXEL](https://www.pexels.com).
Application Asks for your query. enter your desired query to be searched, then asks for maximum number of results.
The `title` and `url` of your request will be stored in `image_result` table. Downloaded images will be in `raw` folder.

### Reliable
This application firstly stores your search query as `InProgress` state, and then tries to process your request, if it wasn't successful, then the status of the request will set to `Failed` otherwise will set to `Successfull`.
The scheduler will search for `Failed` requests, and retry them, and if needed, will update the state of the requests.