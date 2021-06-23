<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::group(['middleware' => 'api'], function () {
    Route::resource(
        '/labels',
        '\App\Http\Controllers\LabelController',
        ['except' => ['update']]
    );

    /*
     * поскольку REST предполагает, что для апдейта существующего ресурса
     * используется PUT, конструкция с ресурными контроллерами ждёт именно PUT
     * однако по ТЗ для обновления должен быть POST, поэтому насильно запрещаем
     * метод в объявлении группы роутов ресурсного контроллера
     * и прописываем его отдельно
     */
    Route::post(
        '/labels/{id}',
        '\App\Http\Controllers\LabelController@update'
    );
});
