# bookmark-maker

This app is built for all those people starting a new job how have no idea what links they will need and for all the ADHD people out there who can't seem to organize their stuff, let alone their browser bookmarks and for anyone else who is just sick and tired of building their own browser bookmarks

This app will take in one or multiple file(s) full of defined links to all the places they want to point at and a file single file that defines the folder structure in the bookmarks in the browser. It will then output an import-able html file for any browser (Chrome and Firefox are tested, open to others if needed).  

You can also set up a file stucture demonstrated under the "[folders](folders/)" directory and generate files using the [BuildFolders workflow](.github/workflows/BuildFolders.yml). This works best if you want many different versions of bookmarks or to make both individual and complex bookmark setups.

#### Setup

1. Create a links.yaml file like [example](example-files/links.yaml)
2. Create a folders.yaml file like [example](example-files/folders.yaml)
3. You can also combine these into one file if you don't care much about the organization
4. Download appropriate bin for your OS (tested on linux and mac), or run in github action (recommended)
5. Run `bookmark-maker build -i example-files/links.yaml -i example-files/folders.yaml -o output.html`

Then in your browser go into the bookmarks manager and import the new "output.html" file.

#### File structure
##### links.yaml
  takes in a list of:
    id: a unique identifier that will get used in pulling folder structure
    name: a user friendly name for the bookmark link
    url: the url for the bookmark link
    
##### folders.yaml
  This is a little more complex, but you can write the folder structure however you want using YAML. You can see a really simple example [here](example-files/folders.yaml). You can find a more complex example with things like name overrides [here](folders/example-1/folders.yaml)
