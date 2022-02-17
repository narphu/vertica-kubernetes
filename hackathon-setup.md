Here are steps to get the MC deployment working for the hackathon:

1. Install go.  Must be version 1.16.2.

curl -LO https://go.dev/dl/go1.16.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz

2. Install docker

3. Install kubectl

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

4. Install helm

curl -LO https://get.helm.sh/helm-v3.8.0-linux-amd64.tar.gz
tar zxvf helm-v3.8.0-linux-amd64.tar.gz
sudo install -o root -g root -m 0755 linux-amd64/helm /usr/local/bin/helm

5. Install kubectx

curl -LO https://github.com/ahmetb/kubectx/releases/download/v0.9.4/kubectx_v0.9.4_linux_x86_64.tar.gz
tar zxvf kubectx_v0.9.4_linux_x86_64.tar.gz
sudo install -o root -g root -m 0755 kubectx /usr/local/bin/kubectx

6. Install kubens

curl -LO https://github.com/ahmetb/kubectx/releases/download/v0.9.4/kubens_v0.9.4_linux_x86_64.tar.gz
tar zxvf kubens_v0.9.4_linux_x86_64.tar.gz
sudo install -o root -g root -m 0755 kubens /usr/local/bin/kubens

7. Install krew

(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

8. Add $HOME/.krew/bin/ to your path in your .bashrc (or whatever login shell you use)

9. Restart shell

10. Change directory in your shell to this git repository

11a. Remove an old kind cluster (if these steps were done before).

scripts.kind.sh term cluster1

11b. Start kubernetes-in-docker (kind).

scripts/kind.sh init cluster1  

If you intend to use nodePort, use the the following command instead:

scripts/kind.sh -p 30001 init cluster1  

12. Setup minio in k8s

scripts/setup-minio.sh -o

13. Install cert-manager in k8s

make install-cert-manager

14. Download a copy of the MC for version 10.0.0.  Copy console rpm to docker-console/packages/vertica-console-x86_64.RHEL6.latest.rpm.

15. Copy MCClient.jar to docker-mcclient/packages and docker-console/packages

16. Build the operator and console containers

make docker-build-operator docker-build-console docker-build-mcclient docker-push-operator docker-push-console docker-push-mcclient

17. Get a copy of a Vertica server RPM and copy it to docker-vertica/packages/vertica-x86_64.RHEL6.latest.rpm

18. Build the vertica container

make docker-build-vertica docker-push-vertica

19. Set this environment variable if you want to use nodePort.

export CONSOLE_NODEPORT=30001

20. Deploy the operator and the console

make deploy

21. Pre-push the vertica server image

scripts/push-to-kind.sh -i vertica/vertica-k8s:latest

22. Create a minIO tenants.  This will serve as the communal path.

kubectl apply -f config/samples/minio.yaml

23. Wait for minIO to come up.  Run this in a loop until you see one of the pods have STATUS column set to completed.

kubectl get pods --selector job-name=create-s3-bucket
NAME                     READY   STATUS      RESTARTS   AGE
create-s3-bucket-4qprk   0/1     Error       0          42s
create-s3-bucket-bbtpp   0/1     Error       0          52s
create-s3-bucket-cvz5q   0/1     Error       0          104s
create-s3-bucket-jppz6   0/1     Completed   0          22s

24. Create the Vertica DB

kubectl apply -f config/samples/v1beta1_verticadb.yaml

25.  Wait for the database to be created

kubectl wait --for=condition=DBInitialized=true vdb/verticadb-sample --timeout 600s


If you are going to use nodePort skip the next 2 steps.

26. If not using node port, expose the MC so you can access it on your localhost

scripts/expose-console.sh

27.  Open up MC in your webbrowser

https://localhost:5450/

28.  The MC database get imported by the operator.  So you need to wait a minute or two after create db for it to be visible in the MC.
