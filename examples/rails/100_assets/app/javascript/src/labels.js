var request = new XMLHttpRequest();
request.open('GET', '/config/env.json', false);  // `false` makes the request synchronous
request.send(null);

if (request.status === 200) {
    document.addEventListener('DOMContentLoaded', function(event) {
        const variables = JSON.parse(request.responseText);
        document.getElementById("environment").innerHTML = variables.environment;
        document.getElementById("link").href = variables.link;

        // Business logic here
        console.log('It works');
        var request_content = new XMLHttpRequest();
        request_content.open('GET', '/api/labels', false);  // `false` makes the request synchronous
        request_content.send(null);
        if (request_content.status === 200) {
            document.getElementById("content").innerHTML = request_content.responseText;
        } else {
            document.getElementById("content").innerHTML = "sorry, error while loading";
        }
    })
}
