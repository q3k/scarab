Scarab, an automation job runner.
=================================

WIP
----
Heavy work in progress. Not ready for use. Don't look at it.

Description
-----------

Scarab aims to help you with the sisyphean task of running automation jobs and keeping note of their progress.

It can be used to replace monolithic CI/CD systems (like Jenkins) with leaner, purpose specific and API-centric job runners. These job runners can each run with their own permission set, they all have a very limited set of functionality (no agent system, no source code control and checkout) and can be either defined using a proto file (text or binary) or embedded/extendedin Go code.

