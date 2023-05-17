
With EFS

```bash
docker run -it --cap-add sys_admin --init public.ecr.aws/amazonlinux/amazonlinux:latest
dnf install -y amazon-efs-utils
mkdir efs
mount -t efs -o tls fs-03a6f6df47f67f7d1:/ efs
```

With NFS

```bash
docker run -it --cap-add sys_admin --init public.ecr.aws/amazonlinux/amazonlinux:latest
dnf install -y nfs-utils
mkdir efs
mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport fs-03a6f6df47f67f7d1.efs.eu-west-1.amazonaws.com:/ efs
```

