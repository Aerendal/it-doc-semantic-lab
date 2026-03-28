# exp_001 — Document Class Detector

## Status
Planned

## Goal
Implement keyword + structure heuristic to detect document class from text.
Input: normalized markdown. Output: class_id (e.g. nlp_algorithm_architecture_spec).

## Input
- corpora/nlp_algorithm_architecture_spec/case_001/normalized/
- schemas/document_classes.yaml

## Expected output
- classifier.py with detect_class() working on case_001
- precision: correct class on all 3 initial corpus cases
- result saved to: experiments/exp_001_document_classification/results/

## Method
1. Load document_classes.yaml
2. For each class: count keyword frequency in heading text + body text
3. Score = sum of keyword hits weighted by heading presence
4. Return class with highest score above threshold

## Acceptance criteria
- detect_class(morfologia_md) == "nlp_algorithm_architecture_spec"
- no false positives on 3 initial cases
