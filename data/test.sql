USE sessionmanager;
BEGIN;
SELECT tap.plan(2);
SELECT tap.has_table(DATABASE(),'users','Check users table');
SELECT tap.has_table(DATABASE(),'user_tokens','Check user_tokens table');
CALL tap.finish();
ROLLBACK;
