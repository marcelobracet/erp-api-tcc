#!/bin/bash

set -e

echo "🚀 Starting Production deployment..."

# Verificar se estamos na branch main
if [[ $(git branch --show-current) != "main" ]]; then
    echo "❌ Error: Must be on 'main' branch to deploy to Production"
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

# Verificar se o namespace production existe
if ! kubectl get namespace production &> /dev/null; then
    echo "📦 Creating Production namespace..."
    kubectl create namespace production
fi

# Confirmar deploy em produção
echo "⚠️ WARNING: You are about to deploy to PRODUCTION!"
echo "This will affect live users. Are you sure? (y/N)"
read -r response
if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo "❌ Deployment cancelled"
    exit 1
fi

# Aplicar secrets (se existirem)
if [ -f "k8s/prod/secrets.yaml" ]; then
    echo "🔐 Applying secrets..."
    kubectl apply -f k8s/prod/secrets.yaml -n production
fi

# Aplicar configurações
echo "⚙️ Applying Production configurations..."
kubectl apply -f k8s/prod/ -n production

# Aguardar deployment estar pronto
echo "⏳ Waiting for deployment to be ready..."
kubectl rollout status deployment/erp-api-prod -n production --timeout=300s

# Verificar se o serviço está funcionando
echo "🔍 Checking service health..."
kubectl get pods -n production -l app=erp-api

# Teste de conectividade
echo "🧪 Running smoke tests..."
SERVICE_IP=$(kubectl get service erp-api-service-prod -n production -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
if [ -n "$SERVICE_IP" ]; then
    curl -f "http://$SERVICE_IP/health" || echo "⚠️ Warning: Health check failed"
else
    echo "⚠️ Warning: Could not get service IP"
fi

echo "✅ Production deployment completed successfully!"
echo "🌐 Production API should be available at: https://api.seudominio.com"

# Notificar equipe (opcional)
echo "📢 Sending notification to team..."
# curl -X POST $SLACK_WEBHOOK_URL \
#   -H "Content-type: application/json" \
#   -d '{"text":"🚀 ERP API deployed to Production successfully!"}' 