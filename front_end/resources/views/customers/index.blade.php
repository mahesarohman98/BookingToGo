@extends('layouts.app')

@section('content')


<h1>Customer List</h1>

<!-- Link to create a new customer -->
<h2><a class="create-link" href="{{ route('customers.create') }}">Create New Customer</a></h2>

<!-- Customer Table -->
<table class="customer-table">
    <thead>
        <tr>
            <th>Name</th>
            <th>Date of Birth</th>
            <th>Phone Number</th>
            <th>Email</th>
            <th>Nationality</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        <!-- Loop through customers and display their information -->
        @foreach ($request['customers'] as $customer)
        <tr>
            <td>{{ $customer['name'] }}</td>
            <td>{{ $customer['dob'] }}</td>
            <td>{{ $customer['phone_number'] }}</td>
            <td>{{ $customer['email'] }}</td>
            <td>{{ $customer['nationality'] }}</td>
            <td>
                <!-- Edit Customer -->
                <a class="edit-link" href="{{ route('customers.edit', ['id' => $customer['id']]) }}">Edit</a>

                <!-- Delete Customer -->
                <form method="post" action="{{ route('customers.destroy', ['id' => $customer['id']]) }}" style="display: inline-block;" onsubmit="return confirm('Are you sure you want to delete this customer?');">
                    @csrf
                    @method('delete')
                    <button type="submit" class="delete-button">Delete</button>
                </form>
            </td>
        </tr>
        @endforeach
    </tbody>
</table>
@endsection