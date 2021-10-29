<?php


namespace App\Http\Controllers;


use App\Helpers\ReportHelper;
use Exception;
use Illuminate\Http\JsonResponse;

/**
 * Class ReportController
 * @package App\Http\Controllers
 */
class ReportController extends Controller
{
    /**
     * @return JsonResponse
     */
    public function generate(): JsonResponse
    {
        $report = ReportHelper::generate();
        if (!$report->saveAsFile()) {
            return $this->respondInternalError('unknown error');
        }
        return $this->respond([
            'result' => true,
        ]);
    }
}
