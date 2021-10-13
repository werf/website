<?php

namespace App\Http\Controllers;

class PingController extends Controller
{
    public function ping(): string
    {
        return 'pong';
    }
}
