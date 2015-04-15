var haiti = (function() {
    var my = {},
        kids = {},
        months = [ "January", "February", "March", "April", "May", "June", "July", "August", "September", "October",
                   "November", "December"];

    function calulateAge(birthdate) {
        var today = new Date();
        var years = today.getFullYear() - birthdate.getFullYear();

        console.log("NOW: " + today);
        console.log("BD:  " + birthdate);

        // Reset birthday to the current year.
        birthdate.setFullYear(today.getFullYear());

        // If the user's birthday has not occurred yet this year, subtract 1.
        if (today < birthdate)
        {
            years--;
        }
        return years;
    }

    my.loadAllChildren = function (viewId) {
        $(".list").html("");
        $(".list").append("<option value=\"new\">Add Child Record</option>");
        $.getJSON('/children?order=FamilyName&order=GivenName', function(result) {
            $.each(result, function(i, field) {
                kids[field.id] = field;
                $(".list").append("<option value=\"" + field.id + "\">" + field.familyname + ", " + field.givenname + "</option>");
            });
            alert(viewId);
            if (viewId != undefined) {
                $(".list").val(viewId);
                my.viewChild(viewId);
            }
        });
    }

    my.clearForm = function() {
        $("#givenname").val("");
        $("#familyname").val("");
        $("#birthdate").val("");
        $("#age").val("");
        $("#grade").val("");
        $("#enteredhousing").val("");
        $("#lefthousing").val("");
        $("#yearshousing").val("");
        $("#mother").val("");
        $("#father").val("");
        $("#siblings").val("");
        $("#village").val("");
        $("#command").html("Save");
    }

    my.viewChild = function (val) {
        if (val == "new") {
            my.newRecord();
            return;
        }
        $("#command").html("New");
        $("#givenname").val(kids[val].givenname);
        $("#familyname").val(kids[val].familyname);
        bd = new Date(kids[val].birthdate);
        console.log(kids[val].birthdate)
        $("#birthdate").val(kids[val].birthdate || "");
        $("#age").val(calulateAge(bd));
        $("#grade").val(kids[val].gradeinschool || "");
        $("#enteredhousing").val(kids[val].enteredhousing || "");
        $("#lefthousing").val(kids[val].lefthousing || "");
        $("#yearshousing").val("");
        $("#mother").val(kids[val].mother || "");
        $("#father").val(kids[val].father || "");
        $("#siblings").val(kids[val].siblings || "");
        $("#village").val(kids[val].village || "");
        console.log(JSON.stringify(kids[val]));
    }

    my.saveRecord = function() {
        // First search based on entered name to see if we have a duplicate
        gn = $("#givenname").val().trim();
        fn = $("#familyname").val().trim();
        $.getJSON('/children?filter=FamilyName=,'+fn+'&filter=GivenName=,'+gn, function(result) {
            if (result.length != 0) {
                alert("Duplicate child name, cannot save.");
                return;
            }
            child = {};
            child.familyname = fn;
            child.givenname = gn;
            alert(JSON.stringify(child));
            $.ajax({
                type: 'POST',
                url: '/children',
                data: JSON.stringify(child),
                success: function(data) {
                    alert(JSON.stringify(data));
                    my.loadAllChildren(data.id);
                },
                contentType: "application/json",
                dataType: 'json',
            });
        });

    }

    my.newRecord = function() {
        $(".list").val("new");
        my.clearForm();
        $(".save").enabled(false);
    }

    my.deleteRecord = function() {
        alert("Are you sure?");
    }

    return my;
}())
