<?php


namespace App\Helpers;


use App\Helpers\ReportHelper\Report;

/**
 * Class ReportHelper
 * @package App\Helpers
 */
class ReportHelper
{
    /**
     * @return Report
     */
    static public function generate(): Report
    {
        $report = new Report();
        $dateTime = now()->format('Y-m-d H:i');
        $report->setContent($dateTime);
        $report->setFilename('report_' . $dateTime . '_');
        $report->setMailSubject('Report@' . $dateTime);
        return $report;
    }
}
