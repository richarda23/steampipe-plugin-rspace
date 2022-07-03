A publicly available Docker image with embedded RSpace dashboard is available.

To run the Dashboard, get the image

    docker pull otter606/rspace-steampipe:0.0.2

Set an API key and URL for your RSpace in environment variables on your host machine, e.g.

    export RSPACE_API_KEY="abcdefg"
    export RSPACE_URL=https://path/to/yourRSpace/

### Basic command to launch dashboard

To run the dashboard:

    docker run --rm  --name rspace-dashboard -p9194:9194 -eRSPACE_API_KEY -eRSPACE_URL otter606/rspace-steampipe:0.0.2

And open your browser at http://localhost:9194/local.dashboard.rspace_dashboard

You should see a dashboard illustrating some stats about RSpace activity. There is an  'Untitled documents' benchmark which lists documents with no name. Having too many untitled documents can make it hard to find or search for content.

![docs/RSpaceDashboard.png](docs/RSpaceDashboard.png)

The Dashboard can show charts, tables and alerts.
The above `docker run` command will remove the Docker container when you stop it using Ctrl-C.

### Running queries

You can run the container as a long-running process:

    docker run -d  --name rspace-dashboard -p9194:9194 -eRSPACE_API_KEY -eRSPACE_URL otter606/rspace-steampipe:0.0.2

and then access the Steampipe query command line to execute arbitrary queries:

    docker exec -it rspace-dashboard steampipe query 

    > select timestamp, payload->>'name' from rspace_event where domain='RECORD' and action='CREATE' and username='bbayham';

### Getting reports in CSV format

You can execute a query on the command line and specify an output format (csv or json):

    docker exec -it rspace-dashboard steampipe query\
     "select timestamp, payload->>'name' as name, payload ->>'id' as id from rspace_event where domain='RECORD' and action='CREATE'" --output csv

### Modifying the dashboard

There are two ways you can access the configuration files that generate the dashboard widgets. Both have identical results. Choose whatever you are most comfortable with. After saving the file, the dashboard automatically updates.

Refer to the [Steampipe docs](https://steampipe.io/docs/mods/writing-dashboards) on how to declare dashboards and controls

#### 1. Modifying the file inside Docker container

You can access the container via

    docker exec -it rspace-dashboard /bin/bash

 vim and nano editors are installed to edit the `control_events.sp` and `rspace_dashboard.sp` files.

#### 2. Mounting the dashboard files directly into the container

Alternatively, you can pull the dashboard project from Github onto your host machine

    git clone --depth 1 -bv0.0.2 https://github.com/richarda23/steampipe-mod-rspace.git

You can now mount the cloned repo inside the Docker container, and edit it with any editor you have installed on your computer. Change the 'src' attribute in the command below to match the location of the cloned repo on your  system.
 
    docker run -d  --name rspace-dashboard -p9194:9194 -eRSPACE_API_KEY -eRSPACE_URL --mount type=bind,src=/absolute/path/to/cloned/gitrepo/steampipe-mod-rspace,dst=/git/mod-rspace otter606/rspace-steampipe:0.0.2