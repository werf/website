<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

/**
 * Class Label
 *
 * @property string $label
 * @package App\Models
 */
class Label extends Model
{
    use HasFactory;

    /**
     * @var bool
     */
    public $timestamps = false;

}
