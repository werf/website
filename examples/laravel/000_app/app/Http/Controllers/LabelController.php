<?php

namespace App\Http\Controllers;

use App\Models\Label;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class LabelController extends Controller
{
    /**
     * @var Label
     */
    private Label $model;

    public function __construct(Label $model)
    {
        $this->model = $model;
    }
    /**
     * Display a listing of the resource.
     *
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        return $this->respond($this->model->all());
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function store(Request $request): JsonResponse
    {
        $label = $request->get('label');
        if (empty($label)) {
            return $this->respondValidationError();
        }
        /** @var Label $resource */
        $resource = (new $this->model);
        $resource->label = $label;
        try {
            $resource->save();
            return $this->respond($resource);
        } catch (Exception $exception) {
            return $this->handleException($exception);
        }
    }

    /**
     * Display the specified resource.
     *
     * @param int $id
     * @return JsonResponse
     */
    public function show(int $id): JsonResponse
    {
        $resource = $this->model->query()->find($id);
        if (!$resource) {
            return $this->respondNotFoundError();
        }
        return $this->respond($this->model->query()->find($id));
    }

    /**
     * Update the specified resource in storage.
     *
     * @param Request $request
     * @param int $id
     * @return JsonResponse
     */
    public function update(Request $request, int $id): JsonResponse
    {
        /** @var Label $resource */
        $resource = $this->model->query()->find($id);
        if (!$resource) {
            return $this->respondNotFoundError();
        }
        $label = $request->get('label');
        if (empty($label)) {
            return $this->respondValidationError();
        }
        $resource->label = $label;
        try {
            $resource->save();
            return $this->respond($resource);
        } catch (Exception $exception) {
            return $this->handleException($exception);
        }
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param int $id
     * @return JsonResponse
     */
    public function destroy(int $id): JsonResponse
    {
        /** @var Label $resource */
        $resource = $this->model->query()->find($id);
        if (!$resource) {
            return $this->respondNotFoundError();
        }
        try {
            $resource->delete();
            return $this->respond([
                'result' => true,
            ]);
        } catch (Exception $exception) {
            return $this->handleException($exception);
        }
    }
}
