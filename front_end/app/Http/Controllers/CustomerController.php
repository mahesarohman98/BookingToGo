<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use GuzzleHttp\Client;

class CustomerController extends Controller
{
    private $httpClient;

    public function __construct()
    {
        $this->httpClient = new Client([
            'base_uri' => "app:3000",
        ]);
    }

    public function index()
    {
        $response = $this->httpClient->get('/customers?per_page=100');
        $request = json_decode($response->getBody(), true);

        return view('customers.index', compact('request'));
    }

    public function create()
    {
        $response = $this->httpClient->get("/countries");
        $countries = json_decode($response->getBody(), true);
        
        return view('customers.create', [
            'countries' => $countries,
        ]);
    }


    public function store(Request $request)
    {

        $request->validate([
            'name' => 'required|string',
            'dob' => 'required|string',
            'phone_number' => 'required|string',
            'email' => 'required|string',
            'country_id' => 'required|integer',  // Validation for integer type
        ]);

        // Extract customer information
        $customerData = $request->only('name', 'dob', 'phone_number', 'email', 'country_id');
        $customerData = [
            'name' => strtolower($request->input('name')),  // Convert name to lowercase
            'dob' => $request->input('dob'),
            'phone_number' => $request->input('phone_number'),
            'email' => $request->input('email'),
            'country_id' => (int) $request->input('country_id'),
        ];


        // Extract family members' information
        $familyMembers = $request->input('familyList', []);

        // Send a POST request to create the customer along with family members
        $response = $this->httpClient->post('/customers', [
            'json' => array_merge($customerData, ['familyList' => $familyMembers]),
        ]);

        // Handle successful creation
        if ($response->getStatusCode() == 201) {
            return redirect()->route('customers.index');
        }

        // Handle errors
        return redirect()->route('customers.create')->with('error', 'Failed to create customer');
    }

    public function edit($id)
    {
        $response = $this->httpClient->get("/countries");
        $countries = json_decode($response->getBody(), true);

        $response = $this->httpClient->get("/customers/{$id}");
        $customer = json_decode($response->getBody(), true);

        return view('customers.edit', [
            'customer' => $customer,
            'countries' => $countries,
        ]);
    }


    public function update(Request $request)
    {

        $id =  $request->input('id');

        $customerData = [
            'name' => strtolower($request->input('name')),  // Convert name to lowercase
            'dob' => $request->input('dob'),
            'phone_number' => $request->input('phone_number'),
            'email' => $request->input('email'),
            'nationality_id' => (int) $request->input('country_id'),
        ];


        // Extract family members' information
        $familyMembers = $request->input('familyList', []);

        // Send a POST request to create the customer along with family members
        $response = $this->httpClient->patch("/customers/{$id}", [
            'json' => array_merge($customerData, ['family_list' => $familyMembers]),
        ]);

        // // Handle successful creation
        if ($response->getStatusCode() == 201) {
            return redirect()->route('customers.index');
        }

        // Handle errors
        return redirect()->route('customers.update')->with('error', 'Failed to update customer');
    }

    public function destroy($id)
    {
        $this->httpClient->delete("/customers/{$id}");
        return redirect()->route('customers.index');
    }

    public function remove_family_member($customer_id, $family_id)
    {
        $response = $this->httpClient->delete("/customers/{$customer_id}/family/{$family_id}");
       
         // Check if the HTTP request was successful
         if ($response->getStatusCode()) {
            return response()->json(['message' => 'Customer created successfully'], 200);
        } else {
            return response()->json([
                'message' => 'Failed to create customer',
                'error' => $response['error'],
            ], $response->getStatusCode());
        }
    }
    
}
