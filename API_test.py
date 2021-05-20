import json
import requests
import xml.etree.ElementTree as ET
import datetime

public_shards = ['admirable', 'animal', 'animaltrait', 'answered', 'banner', 'banners', 'capital', 'category', 'census', 'crime', 'currency', 'customleader', 'customcapital', 'customreligion', 'dbid', 'deaths', 'demonym', 'demonym2', 'demonym2plural', 'dispatches', 'dispatchlist', 'endorsements', 'factbooks', 'factbooklist', 'firstlogin', 'flag', 'founded', 'foundedtime', 'freedom', 'fullname', 'gavote', 'gdp',
                 'govt', 'govtdesc', 'govtpriority', 'happenings', 'income', 'industrydesc', 'influence', 'lastactivity', 'lastlogin', 'leader', 'legislation', 'majorindustry', 'motto', 'name', 'notable', 'policies', 'poorest', 'population', 'publicsector', 'rcensus', 'region', 'religion', 'richest', 'scvote', 'sectors', 'sensibilities', 'tax', 'tgcanrecruit', 'tgcancampaign', 'type', 'wa', 'wabadges', 'wcensus', 'zombie']
private_shards = ['dossier', 'issues', 'issuesummary', 'nextissue',
                  'nextissuetime', 'notices', 'packs', 'ping', 'rdossier', 'unread']

if __name__ == "__main__":
    with open('secrets.json', 'r') as f:
        data = json.load(f)

    headers = {"User-Agent": "PythonAPITest", "X-Password": data["password"]}
    params = {'nation': data["name"], 'q': '+'.join(public_shards)}
    r = requests.get('https://www.nationstates.net/cgi-bin/api.cgi',
                     headers=headers, params=params)
    print(r.text)
    root = ET.fromstring(r.content)
    for child in root:
        print(child.tag, child.attrib, child.text)
        if child.tag == "NEXTISSUETIME":
            next_issue_time = datetime.datetime.fromtimestamp(int(child.text))
            print(next_issue_time)
            time_to_next_issue = next_issue_time - datetime.datetime.now()
            print(time_to_next_issue)
