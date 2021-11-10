<?php

namespace App\Http\Controllers;

use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\Log;

class PingController extends Controller
{
    public function ping(): string
    {
        Log::info('Ping request.');
        return 'pong';
    }
}
