using System;

namespace Jargon
{
    public readonly struct Lemmatizer : Dictionary
    {
        private readonly int _MaxGramLength;
        internal int MaxGramLength => _MaxGramLength;

        private readonly Dictionary _Dictionary;
        internal Dictionary Dictionary => _Dictionary;

        public Lemmatizer(Dictionary d, int maxGramLength)
        {
            if (maxGramLength <= 0) throw new ArgumentException("Must be >= 1", nameof(maxGramLength));

            _Dictionary = d ?? throw new ArgumentNullException(nameof(d));
            _MaxGramLength = maxGramLength;
        }

        public (string Canonical, bool Found) Lookup(string[] s)
        => Dictionary.Lookup(s);
    }
}
