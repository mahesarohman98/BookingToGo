<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\CustomerController;

Route::get('/', function () {
    return redirect()->route('customers.index');
});


Route::get('/customers', [CustomerController::class, 'index'])->name('customers.index');
Route::get('/customers/create', [CustomerController::class, 'create'])->name('customers.create');
Route::post('/customers/create', [CustomerController::class, 'store'])->name('customers.store');
Route::get('/customers/{id}/edit', [CustomerController::class, 'edit'])->name('customers.edit');
Route::post('/customers/edit', [CustomerController::class, 'update'])->name('customers.update');
Route::delete('/customers/{id}/delete', [CustomerController::class, 'destroy'])->name('customers.destroy');

Route::delete('/customers/{customer_id}/family/{family_id}', [CustomerController::class, 'remove_family_member'])->name('customers.remove_family_member');