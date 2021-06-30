Welcome to the worst README.md. I'm working on it ..

goodadvice is a fitness tracking app.

UX is still at MVP level cause I hate UX design



Admins can create daily workouts that are accessible to users and guests.
Those workouts are editable by admins as well.
Users can leave notes about the workout and go back and review them later.
Guests can only view the daily workout. 

Known bugs:
None Currently


Future plans/ideas:
Link user to a Coach
Rename comments table to something more accurate
    Added workout time value to table will need to update all the code for this
Left side menu in framework
    Global setting based on user role
Loved/hated user ratings for workouts on daily workout page for logged in users
    Will need two new tables for this
        One for Loved
        One for Hated
Test/benchmark admin ratings for workouts on Edit workout Page
    Can add column to current table for this
Weekly leaderboard display on daily workout page
    Show name and link to users times for logged in users.
    Only show position and time when not logged in
All time score board link on left menu
    Benchmark workouts
    For each lift and named workout
Addworkout page could autoset the date based on next available date from workout table
Email verification for new users.
User Profiles:
    COMPLETED - Personal Records attached to User accounts
    Set Goals.
Random workout shows up if no workout assigned to the day.
Random workout button.
Search functionality.
Admin page:
    More Roles:
        Coach
User submitted workou

Completed fixes/updates:
V0.0.2
Added
User Profiles
Add and display PRs
Admin functions
Add Movements
Changed Signup and Login display location on index
DDLs are now populated from the Database instead of statically configured
Admin page:
    COMPLETE - Add movements is complete. Admins can now add movements based on sport.
    COMPLETE - Shows Version DB and APP
    COMPLETE - to activate users,
    COMPLETE - to deactivate users,
    COMPLETE - Change roles:
            Admin,
            Moderator, Can be set but doesn't change anything from user access
            User.
COMPLETE -- Saved time separate of user's notes
Validate the input before saving, make sure it's a time stamp
On Workout page display day of week instead of the date
better error handling:
I get kicked out each time it errs,
Cleaning text data before running SQL,
Validate a date on entered workout when saving and give option to fix.
