@extends('layouts.app')

@section('content')


<div class="edit-customer-form">
    <h1>Create Customer</h1>

    <form method="POST" action="{{ route('customers.create') }}">
        @csrf

        <!-- Customer Information -->
        <div>
            <label for="name">Name:</label>
            <input type="text" name="name" required>
        </div>

        <div>
            <label for="dob">Date of Birth: </label>
            <input type="date" name="dob" required>
        </div>

        <div>
            <label for="phone_number">Phone Number:</label>
            <input type="text" name="phone_number" required>
        </div>

        <div>
            <label for="email">Email:</label>
            <input type="email" name="email" required>
        </div>

        <div>
            <label for="country_id">Country ID:</label>
            <select id="country_id" name="country_id" placeholder="Pick a state...">
                @foreach ($countries as $country)
                <option value="{{$country['ID']}}">{{$country['Name']}}</option>
                @endforeach
            </select>
        </div>

        <!-- Family List -->
        <div id="family-list">
            <h2>Family Members</h2>
            <div class="family-member">
                <label for="family_relation[]">Relation:</label>
                <input type="text" name="familyList[0][relation]" required>

                <label for="family_name[]">Name:</label>
                <input type="text" name="familyList[0][name]" required>

                <label for="family_dob[]">Date of Birth:</label>
                <input type="date" name="familyList[0][dob]" required>

                <!-- Delete Family Member Button -->
                <button type="button" class="delete-family-member">Delete</button>
            </div>
        </div>

        <!-- Button to Add More Family Members -->
        <button type="button" id="add-family-member">Add Family Member</button>

        <!-- Submit Form -->
        <div>
            <button type="submit">Create Customer</button>
        </div>
    </form>
</div>

<!-- JavaScript for Adding and Deleting Family Members -->
<script>
    $(document).ready(function() {
        $('select').selectize({
            sortField: 'text'
        });
    });

    document.getElementById('add-family-member').addEventListener('click', () => {
        const familyList = document.getElementById('family-list');
        const familyMemberCount = document.querySelectorAll('.family-member').length;

        const newFamilyMember = `
        <div class="family-member">
            <label for="family_relation[]">Relation:</label>
            <input type="text" name="familyList[${familyMemberCount}][relation]" required>
            
            <label for="family_name[]">Name:</label>
            <input type="text" name="familyList[${familyMemberCount}][name]" required>
            
            <label for="family_dob[]">Date of Birth:</label>
            <input type="date" name="familyList[${familyMemberCount}][dob]" required>

            <button type="button" class="delete-family-member">Delete</button>
        </div>
    `;

        familyList.innerHTML += newFamilyMember;
    });

    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('delete-family-member')) {
            e.target.closest('.family-member').remove();
        }
    });
</script>

@endsection