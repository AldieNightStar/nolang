set text "Create 1000 items list" call @cool-print
call @array1000elems

set text "Multiplied items from first list" call @cool-print
call @multiply-each-elem

set text "Filtered by 5" call @cool-print
set filter @filter-each-5
call @filter-by

set text "Now Filtering by each 50" call @cool-print
set filter @filter-each-50
call @filter-by

set text "Each 100 element is filtered" call @cool-print
set filter @filter-each-100
call @filter-by

ret


:cool-print
	print "========================================="
	print concat "      " text
	print "========================================="
	ret

:fill-array
	set array arr-add array cnt
	ret

:array1000elems
	# Create array
	set array arr-new
	# Fill it with 1000 elements 
	repeat cnt 1000 @fill-array
	# Print out that array
	print array
	# End
	ret

# ==========================================
# ==========================================
# ==========================================

:multiplied-copy
	set array2 arr-add array2 mul item 2
	ret

:multiply-each-elem
	# Create array2 which will have each
	# element of array but multiplied by 2
	set array2 arr-new

	arr-each array item @multiplied-copy

	# Replace old array
	set array array2

	# Print out
	print array

	# End
	ret


# ==========================================
# ==========================================
# ==========================================

# Args
#    IN filter - filter label
#    IN item   - each item of the array (number)
#    OUT array2 - new array
:each-list-filter
	# You see that no "@" at start of filter word
	# It's because label number now in variable
	# So no need to calculate label twice
	call filter
	!if is-ok @each-list-filter-OK
	ret
	:each-list-filter-OK
		set array2 arr-add array2 item
		ret


# Args:
#   IN item - each element number
#   OUT is-ok - is this element ok to add to new list
:filter-each-5
	set is-ok == mod item 5 0
	ret

:filter-each-50
	set is-ok == mod item 50 0
	ret

:filter-each-100
	set is-ok == mod item 100 0
	ret

:filter-by
	# Process each element
	# This time we will specify filter to each parameter

	# Create empty array
	set array2 arr-new

	# Process
	arr-each array item @each-list-filter
	set array array2
	print array
	ret