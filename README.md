# Docker Purge

Docker Purge is a docker tool to remove\clean images, containers, services in batch.

## Usage
#### Removing images
```bash
dkp image -f created>2m3d -f tag=<none>
```
The command above is to remove images that are created **before** 2 months 2 days ago
**and** tagged with `<none>` (not with force).

Available filters are list below.
+ `created`: specifies the create time of an image, in form of `%dy%dm%dd`. e.g. `1y`, `2m`, `3d`, `2m3d`.
+ `name`: specifies the name of an image
+ `tag`: tag of an image
+ `size`: size of an image. e.g. `-f size>=500M`

---
#### Removing containers

```bash
dkp container -f created<3m -f exited>2d 
```
The command above removes containers that are created **after** 3 months ago 
**and** Exited more than 2 days.

Available filter for container
+ `created`: just like images.
+ `exited`: the exited time from now of a container, in form like created.

#### Removing services
**Coming soon**
