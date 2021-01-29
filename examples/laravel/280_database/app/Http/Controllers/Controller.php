<?php

namespace App\Http\Controllers;

use Exception;
use Illuminate\Foundation\Auth\Access\AuthorizesRequests;
use Illuminate\Foundation\Bus\DispatchesJobs;
use Illuminate\Foundation\Validation\ValidatesRequests;
use Illuminate\Http\JsonResponse;
use Illuminate\Routing\Controller as BaseController;

class Controller extends BaseController
{
    use AuthorizesRequests, DispatchesJobs, ValidatesRequests;

    /**
     * @var int
     */
    protected int $statusCode = 200;

    /**
     * @return int
     */
    protected function getStatusCode(): int
    {
        return $this->statusCode;
    }

    /**
     * @param int $statusCode
     * @return $this
     */
    protected function setStatusCode(int $statusCode): static
    {
        $this->statusCode = $statusCode;
        return $this;
    }

    /**
     * @param $data
     * @param array $headers
     * @return JsonResponse
     */
    protected function respond($data, array $headers = []): JsonResponse
    {
        return response()->json($data, $this->getStatusCode(), $headers);
    }

    /**
     * @param string|null $message
     * @return JsonResponse
     */
    protected function respondWithError(?string $message): JsonResponse
    {
        return $this->respond([
            'result' => 'error',
            'comment' => $message,
            'status_code' => $this->getStatusCode()
        ]);
    }

    /**
     * @param string $message
     * @return JsonResponse
     */
    protected function respondNotFoundError($message = 'Объект не найден'): JsonResponse
    {
        return $this->setStatusCode(404)->respondWithError($message);
    }

    /**
     * @param string $message
     * @return JsonResponse
     */
    protected function respondInternalError($message = 'Произошла внутренняя ошибка'): JsonResponse
    {
        return $this->setStatusCode(500)->respondWithError($message);
    }

    /**
     * @param string $message
     * @return JsonResponse
     */
    protected function respondValidationError($message = 'Не все обязательные поля заполнены'): JsonResponse
    {
        return $this->setStatusCode(400)->respondWithError($message);
    }

    /**
     * @param string $message
     * @return JsonResponse
     */
    protected function respondAuthorisationError($message = 'Нет прав'): JsonResponse
    {
        return $this->setStatusCode(401)->respondWithError($message);
    }

    /**
     * @param Exception $e
     * @return JsonResponse
     */
    protected function handleException(Exception $e): JsonResponse
    {
        return $this->respondInternalError($e->getCode() . ': ' . $e->getMessage());
    }
}
