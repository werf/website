<!DOCTYPE html>
<html>
<head>
    <title>werf-guide-app</title>
    <meta name="viewport" content="width=device-width,initial-scale=1">

    <!-- [<ru>] Загрузим CSS-файлы. -->
    <link rel="stylesheet" href="{{ URL::asset('css/site.css') }}" />

    <!-- [<ru>] Загрузим JS-файлы. -->
    <script type="text/javascript" src="{{ URL::asset('js/image.js') }}"></script>
</head>

<body>
<div id="container">
    <!-- [<ru>] Кнопка "Get image", при нажатии на которую будет выполняться Ajax-запрос. -->
    <button type="button" id="show-image-btn">Get image</button>
</div>
</body>
</html>
