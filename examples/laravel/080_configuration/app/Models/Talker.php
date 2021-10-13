<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property string $answer
 * @property string $name
 * @property string $created_at
 * @property string $updated_at
 * @property string $deleted_at
 */
class Talker extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var string[]
     */
    protected $fillable = [
        'answer',
        'name',
    ];
}
