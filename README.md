# About Project

 <p>A functional HTTP/1.1 server from raw TCP sockets in Go without using net/http package. Implemented manual parsing of request line, headers, and body with proper response formatting. Added chunked transfer encoding for big files. </p>
<p>I built this project because I wanted to strengthened my understanding of TCP socket programming and HTTP protocol internals.</p>

<p>Learning resource : <a href="https://www.boot.dev/lessons/b0cebf37-7151-48db-ad8a-0f9399f94c58">https://www.boot.dev/lessons/b0cebf37-7151-48db-ad8a-0f9399f94c58</a></p>

## HTTP 1.1 basic implementation
<br/>

### Example of Chunk Encoding for a big file using this server
<div>
    <p>Reading 32 bytes at a time. (32 in hex is Ox20)</p>
    <img src="./assets/chunked_encoding.png" alt="chunked encoding" />
</div>


### Example of Sending binary data on this server
<div>
    <p>Sending a video using header Content-Type: video/mp4</p>
    <img src="./assets/binary_data.png" alt="binary data" />
</div>



### Example of Serving html on this server
<div>
    <img src="./assets/serving_html.png" alt="serving html curl" />
    <br/>
    <img src="./assets/serving_html_browser.png" alt="serving html browser" />
</div>