<!DOCTYPE html>
<html lang="en">

<head>

    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/selectize.js/0.12.6/js/standalone/selectize.min.js" integrity="sha256-+C0A5Ilqmu4QcSPxrlGpaZxJ04VjsRjKu+G82kl5UJk=" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/selectize.js/0.12.6/css/selectize.bootstrap3.min.css" integrity="sha256-ze/OEYGcFbPRmvCnrSeKbRTtjG4vGLHXgOqsyLFTRjg=" crossorigin="anonymous" />

    <title>Customer Management</title>
</head>

<body>
    <style>
        .edit-customer-form {
            max-width: 600px;
            margin: 0 auto;
            /* Center the form */
            padding: 20px;
            background-color: #f9f9f9;
            border: 1px solid #ddd;
            border-radius: 8px;
        }

        .edit-customer-form h1 {
            text-align: center;
        }

        .edit-customer-form div {
            margin-bottom: 15px;
            /* Space between fields */
        }

        .edit-customer-form label {
            display: block;
            /* Align label with input */
            font-weight: bold;
        }

        .edit-customer-form input,
        .edit-customer-form select {
            width: 100%;
            /* Full-width inputs */
            padding: 8px;
            border: 1px solid #ccc;
            border-radius: 4px;
        }

        .edit-customer-form button {
            background-color: #28a745;
            /* Green for submit button */
            color: white;
            border: none;
            padding: 10px 20px;
            cursor: pointer;
            border-radius: 4px;
        }

        .edit-customer-form button:hover {
            background-color: #218838;
            /* Darker green on hover */
        }

        .delete-family-member {
            background-color: #ff4d4d;
            /* Red for delete */
            color: white;
            border: none;
            padding: 6px 12px;
            cursor: pointer;
        }

        .delete-family-member:hover {
            background-color: #e60000;
            /* Darker red on hover */
        }


        /* Table styling */
        .customer-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }

        .customer-table th,
        .customer-table td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }

        .customer-table th {
            background-color: #f2f2f2;
            font-weight: bold;
        }

        .customer-table tr:hover {
            background-color: #f1f1f1;
            /* Hover effect */
        }

        /* Button styling */
        .delete-button {
            background-color: #ff4d4d;
            /* Red for delete */
            color: white;
            border: none;
            padding: 6px 12px;
            cursor: pointer;
        }

        .delete-button:hover {
            background-color: #e60000;
            /* Darker red on hover */
        }

        .create-link,
        .edit-link {
            color: #007bff;
            /* Bootstrap's blue */
            text-decoration: none;
        }

        .create-link:hover,
        .edit-link:hover {
            text-decoration: underline;
        }
    </style>

    @yield('content')
</body>

</html>