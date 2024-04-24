@extends('layouts.app')

@section('content')
<!-- Style for the form -->

<div class="edit-customer-form">
    <h1>Edit Customer</h1>

    <form method="POST" action="{{ route('customers.update') }}">
        @csrf

        <!-- Customer Information -->
        <div>
            <input type="hidden" name="id" value="{{ $customer['id'] }}" required>
        </div>
        <div>
            <label for="name">Name:</label>
            <input type="text" name="name" value="{{ $customer['name'] }}" required>
        </div>

        <div>
            <label for="dob">Date of Birth:</label>
            <input type="date" name="dob" it value="{{ $customer['dob'] }}" required>
        </div>

        <div>
            <label for="phone_number">Phone Number:</label>
            <input type="text" name="phone_number" value="{{ $customer['phone_number'] }}" required>
        </div>

        <div>
            <label for="email">Email:</label>
            <input type="email" name="email" value="{{ $customer['email'] }}" required>
        </div>

        <div>
            <label for="country_id">Country:</label>
            <select id="country_id" name="country_id">
                @foreach ($countries as $country)
                <option value="{{ $country['ID'] }}" {{ $country['ID'] == $customer['nationality_id'] ? 'selected' : '' }}>
                    {{ $country['Name'] }}
                </option>
                @endforeach
            </select>
        </div>

        <!-- Family List -->
        <div id="family-list">
            <h2>Family Members</h2>
            @foreach ($customer['family_list'] as $i => $family)
            <div class="family-member" data-family-id="{{ $family['id'] }}">
                <input type="hidden" name="familyList[{{ $i }}][id]" value="{{ $family['id'] }}" required>

                <label for="family_relation">Relation:</label>
                <input type="text" name="familyList[{{ $i }}][relation]" value="{{ $family['relation'] }}" required>

                <label for="family_name">Name:</label>
                <input type="text" name="familyList[{{ $i }}][name]" value="{{ $family['name'] }}" required>

                <label for="family_dob">Date of Birth:</label>
                <input type="date" name="familyList[{{ $i }}][dob]" value="{{ $family['dob'] }}" required>

                <!-- Delete Family Member Button -->
                <button type="button"  class="delete-family-member">Delete</button>
            </div>
            @endforeach
        </div>

        <!-- Button to Add More Family Members -->
        <button type="button" id="add-family-member">Add Family Member</button>

        <!-- Submit Form -->
        <button type="submit">Edit Customer</button>
    </form>
</div>

<!-- JavaScript for Adding and Deleting Family Members -->
<script>
    $(document).ready(function() {
        $('select').selectize({
            sortField: 'text'
        });
    });
    document.addEventListener('DOMContentLoaded', function() {
        const customerId = "{{ $customer['id'] }}"; // Customer ID to be used in the request

        document.getElementById('add-family-member').addEventListener('click', () => {
            const familyList = document.getElementById('family-list');
            const familyMemberCount = document.querySelectorAll('.family-member').length;

            const newFamilyMember = `
                <div class="family-member">
                    <label for="family_relation">Relation:</label>
                    <input type="text" name="familyList[${familyMemberCount}][relation]" required>
                    
                    <label for="family_name">Name:</label>
                    <input type="text" name="familyList[${familyMemberCount}][name]" required>
                    
                    <label for="family_dob">Date of Birth:</label>
                    <input type="date" name="familyList[${familyMemberCount}][dob]" required>

                    <button type="button"  class="delete-family-member">Delete</button>
                </div>
            `;

            familyList.innerHTML += newFamilyMember;
        });

        document.addEventListener('click', async (e) => {
            if (e.target.classList.contains('delete-family-member')) {
                const familyMember = e.target.closest('.family-member');
                const familyId = familyMember.getAttribute('data-family-id');

                if (familyId) {
                    try {
                        const response = await fetch(`/customers/${customerId}/family/${familyId}`, {
                            method: 'DELETE',
                            headers: {
                                'Content-Type': 'application/json',
                                'X-CSRF-TOKEN': '{{ csrf_token() }}'
                            }
                        });

                        if (response.ok) {
                            familyMember.remove(); // Remove from the DOM
                        } else {
                            console.error('Error deleting family member');
                        }
                    } catch (error) {
                        console.error('Network error:', error);
                    }
                } else {
                    familyMember.remove(); // If no ID, just remove the DOM element
                }
            }
        });
    });
</script>

@endsection