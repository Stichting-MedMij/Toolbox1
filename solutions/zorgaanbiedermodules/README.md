# Zorgaanbiedermodules

Zorgaanbiedermodules zijn modules die vanuit het PGO gestart kunnen worden. In sommige gevallen wordt de patiënt uitgenodigd door de zorgverlener om een modules uit te voeren, in andere gevallen kan de patiënt op eigen initiatief een zelfhelpmodule uitvoeren

## Zorgaanbiedermodules op uitnodiging van de zorgverlener (usecase 1)

In deze flow gaan we er vanuit dat de patiënt op bezoek is bij de zorgverlener die vervolgens via zijn zorgaanbiederportaal een uitnodiging sturt naar de patiënt om bijvoorbeeld een vragenlijst in te vullen. De bedoeling is dat de patiënt deze module vanuit zijn PGO start, via de DVA doorloopt en dat eventuele resultaten uit de interventie zowel bij de zorgverlener in zijn EPD systeem terug te vinden zijn, als in het PGO van de patiënt.

### Voorstel sequence diagram

```mermaid
sequenceDiagram

autonumber

actor pc as Patiënt op consult
actor zv as Zorgverlener
participant e as EPD
participant zap as Zorgaanbiedersportaal
%% participant Moduleaanbieder as ma

box DVA
    participant dfs as DVA FHIR server
    participant des as DVA email service
    participant dmd as DVA moduledienst
    participant das as DVA auth server
end

participant pgo as PGO

participant mpp as MedMij PGO portaal
actor pt as Patiënt thuis (user agent)

zv ->> e: Open Zorgaanbiedersportaal
e ->> zap: Redirect
zap ->> zv: Toon beschikbare digitale interventies
zv ->> zap: Selecteer passende interventie
zv ->> pc: Vertelt patiënt dat de digitale<br>interventie klaar staat in zijn/haar PGO
zv ->> zap: Nodig patiënt uit voor interventie
zap ->> dfs: POST Task resource voor<br>uitnodiging patiënt met status "requested"
dfs ->> des: Plan email met<br>uitnodiging voor interventie
note over des: Plan in herinneringsemail indien<br>patiënt niet reageert
des ->> pt: Email met link naar keuze PGO's
pt ->> mpp: Open link
mpp ->> pt: Toon PGO's met ondersteuning voor modules
pt ->> mpp: Selecteer gewenste PGO
mpp ->> pt: Redirect naar PGO
note over mpp,pt: Redirect link bevat informatie over<br>zorgaanbieder die interventie heeft gestart
pt ->> pgo: Open PGO
pgo -> pt: Registreer of log in

note over pc,pt:  Verzamelen Task resource met informatie over interventie

pgo ->> pgo: Check voor langdurige toestemming

alt Langdurige toestemming is al gegeven
    pgo ->> das: Haal nieuwe access token op m.b.v. refresh token
    das ->> pgo: Access token + nieuwe refresh token
else Langdurige toestemming is nog niet gegeven
    note over das,pt: Standaard UC Verzamelen authn/authz flow
end

pgo ->> dfs: Haal Task resource op
dfs ->> pgo: Interventie Task resource

pgo ->> pt: Toon pagina met informatie over voorgeschreven interventie

alt Patiënt accepteert interventie
    pt ->> pgo: Accepteer interventie
    pgo ->> dfs: Update Task resource met status "accepted"
else Patiënt accepteert interventie niet
    pt ->> pgo: Weiger interventie
    pgo ->> dfs: Update Task resource met status "rejected"
end


note over pc,pt:  Uitvoeren Zorgaanbiedermodule

pt ->> pgo: Start zorgaanbiedermodule

rect rgb(125,125,125,.2)
  note over pc,pt: Zorgaanbiedermodule authenticatieflow
    pgo ->> pt: Redirect naar DVA auth server
    pt ->> das: Authorization request
    note over das,pt: DigiD
    das ->> pt: Toestemmingsscherm
    note over das,pt: HTTP/1.1 200 Ok<br>Set-Cookie: sessie_cookie
    pt ->> das: Geef toestemming
    das ->> pt: Redirect naar pgo met authz code
    pt ->> pgo: Volg redirect URL met authz code
    pgo ->> das: Access token request
    das ->> pgo: Access token + refresh token
    note over das,pgo: Access token kan later evt.<br>gebruikt worden om resultaat van<br>interventie op te halen
end

pgo ->> dfs: Update Task resource met status "in-progress"
dfs ->> pgo: 200 Ok

pgo ->> pgo: Sla informatie over flow op onder een "state" nonce

pgo ->> pt: Redirect naar DVA moduledienst
pt ->> dmd: Open moduledienst
note over pt,dmd: GET /module?redirect_url=https://pgo.nl/moduleklaar&state=456<br>Host: modules.dva.nl<br>Cookie: sessie_cookie

dmd ->> dmd: Vind sessie informatie m.b.v. sessie cookie
dmd ->> pt: Toon zorgaanbiedermodule
note over dmd,pt: Uitvoeren zorgaanbiedermodule
pt ->> dmd: Zorgaanbiedermodule is afgerond
dmd ->> pt: Redirect terug naar PGO
note over dmd,pt: HTTP/1.1 302 Found<br>Location: https://pgo.nl/moduleklaar?state=456

pt ->> pgo: Open "module afgerond pagina"
note over pgo,pt: GET /moduleklaar?state=456
pgo ->> pgo: Zoek informatie over huidige module flow op m.b.v. "state" parameter

opt Als ophalen van resultaten mogelijk is
    pgo ->> dfs: Haal resultaten van module op
    dfs ->> pgo: Resultaten
end

pgo ->> dfs: Update Task resource met status "completed"
dfs ->> pgo: 200 Ok

pgo ->> pt: Toon module afgerond pagina met resultaten (indien beschikbaar)
```
