<?php


namespace App\Helpers\ReportHelper;


use Illuminate\Mail\Message;
use Illuminate\Support\Facades\Mail;
use Illuminate\Support\Facades\Storage;

class Report
{
    /**
     * @var string
     */
    protected string $content;

    /**
     * @var string
     */
    protected string $filename;

    /**
     * @var string
     */
    protected string $extension = 'txt';

    /**
     * @var string
     */
    protected string $mailSubject;

    /**
     * @return string
     */
    public function getContent(): string
    {
        return $this->content;
    }

    /**
     * @param string $content
     * @return $this
     */
    public function setContent(string $content): static
    {
        $this->content = $content;
        return $this;
    }

    /**
     * @return string
     */
    public function getFilename(): string
    {
        return $this->filename;
    }

    /**
     * @param string $filename
     * @return $this
     */
    public function setFilename(string $filename): static
    {
        $this->filename = $filename;
        return $this;
    }

    /**
     * @return string
     */
    public function getMailSubject(): string
    {
        return $this->mailSubject;
    }

    /**
     * @param string $mailSubject
     * @return $this
     */
    public function setMailSubject(string $mailSubject): static
    {
        $this->mailSubject = $mailSubject;
        return $this;
    }

    /**
     * @return bool
     */
    public function saveAsFile(): bool
    {
        $filename = $this->filename . '.' . $this->extension;
        return Storage::disk('s3')->put($filename, $this->content);
    }

    /**
     * @param string $recipient
     * @return void
     */
    public function sendByMail(string $recipient): void
    {
        Mail::raw($this->content, function (Message $message) use ($recipient) {
            $message->to($recipient)
                ->subject($this->mailSubject);
        });
    }
}
