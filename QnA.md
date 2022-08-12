How long did this assignment take?
> I had to come up to speed on Go, but outside of that, about 5h.

What was the hardest part?
> Structuring the project such that the dependencies played well together. Also,
> the ORM's API is a bit unintuitive; I don't love passing a pointer argument
> and having it's data change.

Did you learn anything new?
> I did. As a relative newbie to Go, pretty much everything was newly-learned. I
> especially liked the `gorilla` toolkit for building web servers.

Is there anything you would have liked to implement but didn't have the time to?
> Because the requirements didn't state a need for security, the authentication
> mechanism is intentionally naïve. Storing passwords in plain text just _feels_
> gross, so the first change I would make is encrypting those (hence the use of
> the `unencrypted_password`/`UnencryptedPassword` column/attribute. A more
> robust authorization toolkit would probably be next.

What are the security holes (if any) in your system? If there are any, how would you fix them?
> See above regarding unencrypted passwords. Huge, glaring security hole.
> Relatively easily fixed tho, by simply adding a salt and encryption routine,
> and storing the encrypted result instead.
>
> In addition to that, the JWT-based security as implemented is uncomfortably loose.  
> A given email can have multiple tokens provisioned (different but equally-valid
> tokens are returned by the calls to `POST /signups` and `POST /logins`). On
> top of that, there's no security check with respect to the JWT payload itself.
> Simply checking the `aud` and `iss` claims would be an improvement, as would
> validating the signature. Lastly, using an asymmetric encryption algorithm
> would be a must-have for a real-world application.
>
> In general, I'm not a fan of rolling one's own auth system; it's too easy to
> get it wrong and doing so can have disasterous consequences. I would prefer to
> use a vendor (e.g. Auth0, etc.), but that would have substantically increased
> complexity for the given requirements.

Do you feel that your skills were well tested?
> Yes.
