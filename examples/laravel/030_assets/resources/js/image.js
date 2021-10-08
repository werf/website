window.onload=function(){
  var btn = document.getElementById("show-image-btn")
  btn.addEventListener("click", _ => {
      fetch("/images/werf-logo.svg")
        .then((data) => data.text())
        .then((html) => {
          const svgContainer = document.getElementById("container")
          svgContainer.insertAdjacentHTML("beforeend", html)
          var svg = svgContainer.getElementsByTagName("svg")[0]
          svg.setAttribute("id", "image")
          btn.remove()
        }
      )
    }
  )
}

