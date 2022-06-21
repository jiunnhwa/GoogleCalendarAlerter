# GoogleCalendarAlerter
A console golang program to retrieve upcoming events from google calendar, and sounds an alarm on event time.
Uses OAuth2, it will get authorization token from web, if none exists, and saves the refresh token on file to renew expiration.
The service will fetch at the start of each minute.


## Screen Capture - Retrieves events up till tomorrow every minute. ##
![Fetch Events](https://github.com/jiunnhwa/GoogleCalendarAlerter/blob/main/20220621%20MyGoogleCalAlerter.gif?raw=true "Send OTP")



## Research Notes: ##
- [Google Calendar for Developers - Calendar API](https://developers.google.com/calendar/api/quickstart/go) - Google's sample to create a simple Go command-line application that makes requests to the Google Calendar API. 
- [Develop Google Calendar solutions](https://developers.google.com/calendar) - Overview of enhancing and automating the Google Calendar experience.
- [OAuth Playground](https://developers.google.com/oauthplayground/)
- [An Introduction to OAuth 2](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2) - An overview of OAuth 2 roles, authorization grant types, use cases, and flows.
- [Using OAuth 2.0 for Web Server Applications](https://developers.google.com/identity/protocols/oauth2/web-server) - How to configure OAuth.
- [Send GET Request with HTTP / Webhook APIon New or Updated Event (Instant) from Google Calendar API](https://pipedream.com/apps/google-calendar/integrations/http/get-request-with-http-webhook-api-on-new-or-updated-event-instant-from-google-calendar-api-int_EzsYAz) - 3rd party solution(pipedream).
- [How to Integrate The Google Calendar API Into Your Web Application](https://qalbit.com/blog/how-to-integrate-the-google-calendar-api-into-your-web-application/) - 3rd party solution(qalbit).
- [Google Calendar Webhooks with Node.js](https://fusebit.io/blog/google-calendar-webhooks)  - 3rd party solution(fusebit).
- [How to Use the Google Calendar API](https://zapier.com/engineering/how-to-use-the-google-calendar-api/)  - 3rd party solution(zapier).
- [How to Integrate The Google Calendar API Into Your App](https://www.nylas.com/blog/integrate-google-calendar-api)  - 3rd party solution(nylas).
- [Creating a Google API Key](https://docs.simplecalendar.io/google-api-key/)  - 3rd party solution(simplecalendar).
- [FullCalendar can display events from a public Google Calendar](https://fullcalendar.io/docs/google-calendar) - 3rd party solution(fullcalendar).
- [Google Calendar Simple API](https://github.com/kuzmoyev/google-calendar-simple-api) - Pythonic wrapper for the Google Calendar API.
- [Parsing Google Calendar events with Python](https://qxf2.com/blog/google-calendar-python/) - Pythonic sample to look for all-day events with the word ‘PTO’ in the event title.
- [Google Calendar API with Python](https://dev.to/nelsoncode/google-calendar-api-con-python-1ib1) - Pythonic sample of the main methods to interact with Google Calendar API.
- [How to Integrate and Synchronize Google Calendar with Your Blazor Application](https://www.grapecity.com/blogs/how-to-integrate-synchronize-google-calendar-with-blazor-application) - Blazor C# sample using ComponentOne.
- [Google calendar API integration in .Net Core](https://www.thecodehubs.com/google-calendar-api-integration-in-net-core/) - c# .Net Core MVC Web Application integrating google calendar.
- [Google Calendar API Authentication with C#](https://www.daimto.com/google-calendar-api-authentication-with-c/) - how to access the Google Calendar API with an authenticated user. 
- [How to get a Google access token with CURL](https://www.daimto.com/how-to-get-a-google-access-token-with-curl/) - In order to tell the authorization server that we want the token returned in the web browser we just send it urn:ietf:wg:oauth:2.0:oob as our redirect uri.
- [Calendar API Samples for .NET](https://github.com/LindaLawton/Google-Dotnet-Samples/tree/master/Samples/Calendar%20API) - These samples show how to access the Calendar API with the Offical Google .Net client library by Linda Lawton.
