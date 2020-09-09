var request = new XMLHttpRequest();
request.open('GET', '/config/env.json', false);  // `false` makes the request synchronous
request.send(null);

if (request.status === 200) {
    const variables = JSON.parse(request.responseText);
    document.getElementById("url").href = variables.url;

    // Business logic here
    console.log('It works');
    var request_content = new XMLHttpRequest();
    request_content.open('GET', '/api/', false);  // `false` makes the request synchronous
    request_content.send(null);
    if (request_content.status === 200) {
        document.getElementById("content").innerHTML = request_content.responseText;
    } else {
        document.getElementById("content").innerHTML = "sorry, error while loading";
    }

}