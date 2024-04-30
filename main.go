package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)


func main() {
	var err error
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/", View)

	server := &http.Server{
        Addr: ":8080",
    }
	go func() {
		err = server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
    }()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	<-appSignal

	server.Close()
	
}

func View(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://localhost:3333/visit")
	if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	data := WebLogResponse{}
	json.Unmarshal(responseData, &data)
	io.WriteString(w, fmt.Sprintf(Html, data.Count))
}

type WebLogResponse struct {
	Count int `json:"count"`
}

var Html string = `
<!DOCTYPE html>
<html>
<title>W3.CSS Template</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Montserrat">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
<style>
body, h1,h2,h3,h4,h5,h6 {font-family: "Montserrat", sans-serif}
.w3-row-padding img {margin-bottom: 12px}
/* Set the width of the sidebar to 120px */
.w3-sidebar {width: 120px;background: #222;}
/* Add a left margin to the "page content" that matches the width of the sidebar (120px) */
#main {margin-left: 120px}
/* Remove margins from "page content" on small screens */
@media only screen and (max-width: 600px) {#main {margin-left: 0}}
</style>
<body class="w3-black">



<!-- Page Content -->
<div class="w3-padding-large" id="main">
  <!-- Header/Home -->
  <header class="w3-container w3-padding-32 w3-center w3-black" id="home">
    <h1 class="w3-jumbo"><span class="w3-hide-small">I'm</span> Husni Faiz.</h1>
    <p>Software Engineer</p>

	<p> %d website visits!</p>
  </header>

  <!-- About Section -->
  <div class="w3-content w3-justify w3-text-grey w3-padding-64" id="about">
    <h2 class="w3-text-light-grey">About</h2>
    <hr style="width:200px" class="w3-opacity">
    <p>I currently work as a back-end software engineer. I work with Go, Kafka and more. Previously worked with C, Real-time OS, Bootloaders etc.
    </p>
    <h3 class="w3-padding-16 w3-text-light-grey">My Skills (Not Verified)</h3>
    <p class="w3-wide">Engineering</p>
    <div class="w3-white">
      <div class="w3-dark-grey" style="height:28px;width:95%%"></div>
    </div>
    <p class="w3-wide">Debugging</p>
    <div class="w3-white">
      <div class="w3-dark-grey" style="height:28px;width:98%%"></div>
    </div>
<!-- 
    <button class="w3-button w3-light-grey w3-padding-large w3-section">
      <i class="fa fa-download"></i> Download Resume
    </button> -->

    <!-- <h3 class="w3-text-light-grey">Social Platforms</h2> -->
    <div class="w3-content w3-padding-64 w3-text-grey w3-xlarge">
        <!-- <h3>Social Platforms</h3> -->
        <a href="https://twitter.com/husni__faiz" target="_blank" class="w3-hover-text-green"><i class="fa fa-twitter w3-hover-opacity"></i></a>
        <!-- <a href="https://www.linkedin.com/in/husni-faiz/" target="_blank" class="w3-hover-text-green"><i class="fa fa-linkedin w3-hover-opacity"></i></a> -->
    </div>
 
  <!-- End About Section -->
  </div>
  

<!-- END PAGE CONTENT -->
</div>

</body>
</html>
`