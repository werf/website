<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;

class S3FileController extends Controller
{
    private $fileName = "filename";

    public function download()
    {
        if (Storage::disk('s3')->missing($this->fileName)) {
            return "You haven't uploaded anything yet.";

        }

        return Storage::disk('s3')->download($this->fileName);
    }

    public function upload(Request $request): string
    {
        if (!$request->hasFile('file')) {
            return "You didn't pass the file to upload :(";
        }

        Storage::disk('s3')->putFileAs(
            '', $request->file('file'), $this->fileName
        );

        return "File uploaded.";
    }
}
