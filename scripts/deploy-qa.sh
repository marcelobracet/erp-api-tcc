#!/bin/bash

set -e

echo "🚀 Starting QA deployment..."

# Verificar se estamos na branch develop
if [[ $(git branch --show-current) != "develop" ]]; then
    echo "❌ Error: Must be on 'develop' branch to deploy to QA"
    exit 1
fi

# Verificar se há mudanças não commitadas
if [[ -n $(git status --porcelain) ]]; then
    echo "❌ Error: There are uncommitted changes. Please commit or stash them first."
    exit 1
fi

# Verificar se o kubectl está configurado
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Error: kubectl is not configured or cluster is not accessible"
    exit 1
fi

# Verificar se o namespace qa existe
if ! kubectl get namespace qa &> /dev/null; then
    echo "📦 Creating QA namespace..."
    kubectl create namespace qa
fi

# Aplicar secrets (se existirem)
if [ -f "k8s/qa/secrets.yaml" ]; then
    echo "🔐 Applying secrets..."
    kubectl apply -f k8s/qa/secrets.yaml -n qa
fi

# Aplicar configurações
echo "⚙️ Applying QA configurations..."
kubectl apply -f k8s/qa/ -n qa

# Aguardar deployment estar pronto
echo "⏳ Waiting for deployment to be ready..."
kubectl rollout status deployment/erp-api-qa -n qa --timeout=300s

# Verificar se o serviço está funcionando
echo "🔍 Checking service health..."
kubectl get pods -n qa -l app=erp-api

# Teste de conectividade
echo "🧪 Running smoke tests..."
SERVICE_IP=$(kubectl get service erp-api-service-qa -n qa -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
if [ -n "$SERVICE_IP" ]; then
    curl -f "http://$SERVICE_IP/health" || echo "⚠️ Warning: Health check failed"
else
    echo "⚠️ Warning: Could not get service IP"
fi

echo "✅ QA deployment completed successfully!"
echo "🌐 QA API should be available at: https://qa-api.seudominio.com" 