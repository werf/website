<?php

namespace App\Http\Controllers;

use App\Models\Talker;
use Illuminate\Http\Request;
use Illuminate\Routing\Controller;

class TalkerController extends Controller
{
    public function say(): string
    {
        /** @var Talker $talker */
        $talker = Talker::first();

        if (!$talker) {
            return "I have nothing to say";
        }

        return $talker->answer . ", " . $talker->name . "!";
    }

    public function remember(Request $request): string
    {
        if (!$request->has('answer')) {
            return "You forgot the answer :(";
        }

        if (!$request->has('name')) {
            return "You forgot the name :(";
        }

        /** @var Talker $talker */
        $talker = Talker::firstOrNew();
        $talker->answer = $request->input('answer');
        $talker->name = $request->input('name');
        $talker->save();

        return "Got it.";
    }
}
