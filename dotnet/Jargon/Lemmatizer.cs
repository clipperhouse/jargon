using System;
using System.Collections.Generic;
using System.IO;

namespace Jargon
{
    public readonly struct Lemmatizer : ILemmatizingDictionary
    {
        private readonly int _MaxGramLength;
        internal int MaxGramLength => _MaxGramLength;

        private readonly ILemmatizingDictionary _Dictionary;
        internal ILemmatizingDictionary Dictionary => _Dictionary;

        public Lemmatizer(ILemmatizingDictionary d, int maxGramLength)
        {
            if (maxGramLength <= 0) throw new ArgumentException("Must be >= 1", nameof(maxGramLength));

            _Dictionary = d ?? throw new ArgumentNullException(nameof(d));
            _MaxGramLength = maxGramLength;
        }

        public (string Canonical, bool Found) Lookup(string[] s)
        => Dictionary.Lookup(s);

        public IEnumerable<Token> Lemmatize(string str) => Jargon.Lemmatize(str, this);
        public IEnumerable<Token> Lemmatize(TextReader reader) => Jargon.Lemmatize(reader, this);
        public IEnumerable<Token> Lemmatize(IEnumerable<Token> tokens) => Jargon.Lemmatize(tokens, this);
    }
}
