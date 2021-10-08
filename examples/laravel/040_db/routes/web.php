<?php

use App\Http\Controllers\PingController;
use App\Http\Controllers\TalkerController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::get('/', function () {
    return view('welcome');
});

Route::get('/ping', [PingController::class, 'ping']);

Route::get('/image', function () {
    return view('image');
});

Route::get('/say', [TalkerController::class, 'say']);
Route::get('/remember', [TalkerController::class, 'remember']);
